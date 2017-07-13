package lua

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Laughs-In-Flowers/shiva/lib/xrror"
	l "github.com/yuin/gopher-lua"
)

type Lua struct {
	Configuration
	*l.LState
}

func New(debug bool, cnf ...Config) (*Lua, error) {
	L := &Lua{
		nil,
		l.NewState(l.Options{
			SkipOpenLibs:        true,
			IncludeGoStackTrace: debug,
		}),
	}
	L.Configuration = newConfiguration(L, cnf...)
	err := L.Configure()
	if err != nil {
		return nil, err
	}
	return L, nil
}

type ConfigFn func(*Lua) error

type Config interface {
	Order() int
	Configure(*Lua) error
}

type config struct {
	order int
	fn    ConfigFn
}

func DefaultConfig(fn ConfigFn) Config {
	return config{50, fn}
}

func NewConfig(order int, fn ConfigFn) Config {
	return config{order, fn}
}

func (c config) Order() int {
	return c.order
}

func (c config) Configure(L *Lua) error {
	return c.fn(L)
}

type configList []Config

func (c configList) Len() int {
	return len(c)
}

func (c configList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c configList) Less(i, j int) bool {
	return c[i].Order() < c[j].Order()
}

type Configuration interface {
	Add(...Config)
	AddFn(...ConfigFn)
	Configure() error
	Configured() bool
}

type configuration struct {
	L          *Lua
	configured bool
	list       configList
}

func newConfiguration(L *Lua, conf ...Config) *configuration {
	c := &configuration{
		L:    L,
		list: builtIns,
	}
	c.Add(conf...)
	return c
}

func (c *configuration) Add(conf ...Config) {
	c.list = append(c.list, conf...)
}

func (c *configuration) AddFn(fns ...ConfigFn) {
	for _, fn := range fns {
		c.list = append(c.list, DefaultConfig(fn))
	}
}

func configure(L *Lua, conf ...Config) error {
	for _, c := range conf {
		err := c.Configure(L)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *configuration) Configure() error {
	sort.Sort(c.list)

	err := configure(c.L, c.list...)
	if err == nil {
		c.configured = true
	}

	return err
}

func (c *configuration) Configured() bool {
	return c.configured
}

var builtIns = []Config{
	config{1001, lRegistryIndex},
	//lconfig{1002, lReservedReferences},
	config{1003, lRequireFunc},
	config{1004, lTracebackFunc},
	config{1005, lModules},
	config{1006, lDetails},
}

var registryIndex *l.LTable

func lRegistryIndex(L *Lua) error {
	registryIndex = L.Get(l.RegistryIndex).(*l.LTable)
	return nil
}

func (L *Lua) RegistryIndex() *l.LTable {
	return registryIndex
}

//func lReservedReferences(L *Lua) error {
//	for _, v := range reservedReferences {
//		registryIndex.RawSetInt(v, l.LBool(true))
//	}
//	return nil
//}

//const RESERVED_REFS_START = -20000

//const (
//	RESERVED_REFS_HEAD = iota + RESERVED_REFS_START
//	LUA_TRACEBACK_FUNC
//	RESERVED_REFS_TAIL
//)

//var reservedReferences = []int{
//	LUA_TRACEBACK_FUNC,
//	RESERVED_REFS_TAIL,
//}

var requireFn *l.LFunction

func require(L *l.LState) int {
	var loopdetection = &l.LUserData{}
	name := L.CheckString(1)
	loaded := L.GetField(registryIndex, "_LOADED")
	lv := L.GetField(loaded, name)
	if l.LVAsBool(lv) {
		if lv == loopdetection {
			L.RaiseError("loop or previous error loading module: %s", name)
		}
		L.Push(lv)
		return 1
	}
	loaders, ok := L.GetField(registryIndex, "_LOADERS").(*l.LTable)
	if !ok {
		L.RaiseError("package.loaders must be a table")
	}
	messages := []string{}
	var modasfunc l.LValue
	for i := 1; ; i++ {
		loader := L.RawGetInt(loaders, i)
		if loader == l.LNil {
			L.RaiseError("module %s not found:\n\t%s, ", name, strings.Join(messages, "\n\t"))
		}
		L.Push(loader)
		L.Push(l.LString(name))
		L.Call(1, 1)
		ret := L.Get(L.GetTop())
		switch retv := ret.(type) {
		case *l.LFunction:
			modasfunc = retv
			goto loopbreak
		case l.LString:
			messages = append(messages, string(retv))
		}
	}
loopbreak:
	L.SetField(loaded, name, loopdetection)
	L.Push(modasfunc)
	L.Push(l.LString(name))
	L.Call(1, 1)
	ret := L.Get(L.GetTop())
	modv := L.GetField(loaded, name)
	if ret != l.LNil && modv == loopdetection {
		L.SetField(loaded, name, ret)
		L.Push(ret)
	} else if modv == loopdetection {
		L.SetField(loaded, name, l.LTrue)
		L.Push(l.LTrue)
	} else {
		L.Push(modv)
	}
	return 1
}

func lRequireFunc(L *Lua) error {
	requireFn = L.NewClosure(require)
	L.SetGlobal("require", requireFn)
	return nil
}

func (L *Lua) RequireFn() *l.LFunction {
	return requireFn
}

var tracebackFn *l.LFunction

func traceback(L *Lua) func(*l.LState) int {
	return func(*l.LState) int {
		L.CallShv("_traceback", 1, 1)
		return 1
	}
}

func lTracebackFunc(L *Lua) error {
	tracebackFn = L.NewClosure(traceback(L))
	registryIndex.RawSetString("_traceback", tracebackFn)
	return nil
}

func (L *Lua) TracebackFn() *l.LFunction {
	return tracebackFn
}

var modulesToAdd []GetModule = defaultModules

func lModules(L *Lua) error {
	var err error
	for _, fn := range modulesToAdd {
		m := fn(L)
		if err = L.CallByParam(l.P{
			Fn:      L.NewFunction(m.Loader),
			NRet:    0,
			Protect: true,
		}, l.LString(m.Tag())); err != nil {
			return err
		}
		err = m.LoadEmb(L)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetModules(g ...GetModule) Config {
	return NewConfig(50,
		func(L *Lua) error {
			for _, gm := range g {
				modulesToAdd = append(modulesToAdd, gm)
			}
			return nil
		})
}

func SetPath(k, v string) Config {
	return NewConfig(9000,
		func(L *Lua) error {
			luaPath(L, k, v)
			return nil
		})
}

func luaPath(L *Lua, k string, v string) {
	c := L.GetGlobal(k)
	ls := l.LString(v)
	if c != ls {
		setPathWithLua(L, v)
		L.SetGlobal(k, ls)
	}
}

func setPathWithLua(L *Lua, path string) {
	tb := L.GetField(L.Get(l.EnvironIndex), "package").(*l.LTable)
	np := fmt.Sprintf("%s/?.lua;%s/?;/usr/local/share/lua/5.1/?.lua;/usr/local/share/lua/5.1/?/init.lua", path, path)
	tb.RawSetString("path", l.LString(np))
}

type Detail struct {
	Key, Value string
}

var (
	luaVers = &Detail{"version", "TBD"}
	luaPlat = &Detail{"platform", PLATFORM}
	luaArch = &Detail{"architechture", ARCH}

	luaDetails = []*Detail{
		luaVers, luaPlat, luaArch,
	}
)

func setDetails(L *Lua, ds ...*Detail) {
	for _, d := range ds {
		L.SetField(lShvModule, d.Key, l.LString(d.Value))
	}
}

func lDetails(L *Lua) error {
	setDetails(L, luaDetails...)
	return nil
}

func SetDetails(version, platform, arch string) Config {
	return NewConfig(50,
		func(L *Lua) error {
			luaVers.Value = version
			luaPlat.Value = platform
			luaArch.Value = arch
			return nil
		})
}

func (L *Lua) CallWithTraceback(nargs, nresults int) error {
	err := L.PCall(nargs, nresults, tracebackFn)
	if err != nil {
		return err
	}
	return nil
}

func (L *Lua) PushString(s string) {
	L.Push(l.LString(s))
}

func (L *Lua) CheckNargs(n int) {
	nargs := L.GetTop()
	if nargs < n {
		if n == 1 {
			L.RaiseError("expecting at least 1 argument")
		} else {
			L.RaiseError("expecting at least %d arguments", n)
		}
	}
}

var NoTableError = xrror.Xrror("no table named %s").Out

type TakeUserData func(*l.LUserData)

func PushNewUserData(L *l.LState, fn TakeUserData, mt string) {
	ud := L.NewUserData()
	fn(ud)
	L.SetMetatable(ud, L.GetTypeMetatable(mt))
	L.Push(ud)
}

func (L *Lua) GetModule(name string) (*l.LTable, error) {
	g := L.FindTable(L.Get(l.RegistryIndex).(*l.LTable), "_LOADED", 1)
	tb := L.GetField(g, name).(*l.LTable)
	t := tb.Type()
	if t != l.LTNil && t == l.LTTable {
		return tb, nil
	}
	return nil, NoTableError(name)
}

func (L *Lua) RegisterFuncsOn(name string, funcs map[string]l.LGFunction) {
	if tb, err := L.GetModule(name); err == nil {
		L.SetFuncs(tb, funcs)
	}
}

func NewProperty(get, set l.LGFunction) l.LGFunction {
	return func(L *l.LState) int {
		switch L.GetTop() {
		case 2:
			if get != nil {
				return get(L)
			}
			L.RaiseError("cannot get %s property", L.Get(2))
		case 3:
			if set != nil {
				return set(L)
			}
			L.RaiseError("cannot set %s property", L.Get(2))
		default:
			L.RaiseError("error accessing property %s", L.Get(L.GetTop()))
		}
		return 0
	}
}

//func initEmbedded(L *Lua) error {
//var err error
//for _, m := range L.modules {
//	err = m.LoadEmb(L)
//	if err != nil {
//		return err
//	}
//}
//return err
//}

//func dbgLua(L *Lua, e *Engine) error {
//if e.debug {
// run debug operations
//}
//return nil
//}

//func tableSpew(t *l.LTable) {
//t.ForEach(func(k, v l.LValue) {
//	spew.Dump(k, v)
//})
//}

//var ShvMetaTables *Tables

//func init() {
//	ShvMetaTables = NewTables()
//}

//SHV_PARAM_NAME_STRING_TABLE = iota + SHV_RESERVED_REFS_START
//SHV_TAG_TABLE
//SHV_WINDOW_TABLE
//SHV_MODULE_TABLE
//SHV_ACTION_TABLE
//SHV_METATABLE_REGISTRY
//SHV_ROOT_AUDIO_NODE
//MT_shv_window
//MT_shv_program
//MT_shv_texture2d
//MT_shv_framebuffer
//MT_shv_image_buffer
//MT_shv_scene_node
//MT_shv_wrap_node
//MT_shv_program_node
//MT_shv_bind_node
//MT_shv_read_uniform_node
//MT_shv_translate_node
//MT_shv_scale_node
//MT_shv_rotate_node
//MT_shv_lookat_node
//MT_shv_billboard_node
//MT_shv_blend_node
//MT_shv_depth_test_node
//MT_shv_viewport_node
//MT_shv_color_mask_node
//MT_shv_cull_face_node
//MT_shv_cull_sphere_node
//MT_shv_draw_node
//MT_shv_pass_filter_node
//MT_shv_audio_buffer
//MT_shv_audio_node
//MT_shv_gain_node
//MT_shv_lowpass_filter_node
//MT_shv_highpass_filter_node
//MT_shv_audio_track_node
//MT_shv_audio_stream_node
//MT_shv_oscillator_node
//MT_shv_spectrum_node
//MT_shv_capture_node
//MT_shv_buffer
//MT_shv_buffer_view
// The following list must be kept in sync with shv_buffer_view_type
//MT_VIEW_TYPE_float
//MT_VIEW_TYPE_vec2
//MT_VIEW_TYPE_vec3
//MT_VIEW_TYPE_vec4
//MT_VIEW_TYPE_ubyte
//MT_VIEW_TYPE_byte
//MT_VIEW_TYPE_ubyte_norm
//MT_VIEW_TYPE_byte_norm
//MT_VIEW_TYPE_ushort
//MT_VIEW_TYPE_short
//MT_VIEW_TYPE_ushort_elem
//MT_VIEW_TYPE_ushort_norm
//MT_VIEW_TYPE_short_norm
//MT_VIEW_TYPE_uint
//MT_VIEW_TYPE_int
//MT_VIEW_TYPE_uint_elem
//MT_VIEW_TYPE_uint_norm
//MT_VIEW_TYPE_int_norm
//MT_VIEW_TYPE_END_MARKER
//MT_shv_vec2
//MT_shv_vec3
//MT_shv_vec4
//MT_shv_mat2
//MT_shv_mat3
//MT_shv_mat4
//MT_shv_quat
//MT_shv_http_request
//MT_shv_socket
//MT_shv_rand
//ENUM_shv_buffer_view_type,
//ENUM_shv_texture_format,
//ENUM_shv_texture_type,
//ENUM_shv_texture_min_filter,
//ENUM_shv_texture_mag_filter,
//ENUM_shv_texture_wrap,
//ENUM_shv_depth_func,
//ENUM_shv_cull_face_mode,
//ENUM_shv_draw_mode,
//ENUM_shv_blend_mode,
//ENUM_shv_window_mode,
//ENUM_shv_display_orientation,
