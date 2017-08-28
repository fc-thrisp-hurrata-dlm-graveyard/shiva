package engine

import (
	"os"
	"sort"
	"time"

	"github.com/Laughs-In-Flowers/log"
	"github.com/Laughs-In-Flowers/shiva/lib/display"
	"github.com/Laughs-In-Flowers/shiva/lib/ecs"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/input"
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/Laughs-In-Flowers/shiva/lib/math"
	"github.com/Laughs-In-Flowers/shiva/lib/render"
	"github.com/Laughs-In-Flowers/shiva/lib/scene"
	"github.com/Laughs-In-Flowers/shiva/lib/xrror"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type HandleErrorFunc func(*Engine, error)

type LoadLuaModuleFunc func() error

type Engine struct {
	Configuration
	log.Logger
	hefn    HandleErrorFunc
	llfn    LoadLuaModuleFunc
	last    error
	restart bool
	kill    bool
	debug   bool
	FPS     int64
	gp      graphics.Provider
	rs      string
	l       *lua.Lua
	w       *ecs.World
}

func New(debug bool, cnf ...Config) (*Engine, error) {
	e := &Engine{debug: debug}
	c := newConfiguration(e, cnf...)
	e.Configuration = c
	err := e.Configure()
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Engine) Run() {
	handleError := e.hefn
	loadLua := e.llfn
	world := e.w
	e.Print("running...")
RESTART:
	reload(e, handleError, loadLua)
	for {
		switch {
		case e.restart:
			e.restart = false
			e.Print("restarting...")
			goto RESTART
		case e.kill:
			goto QUIT
		default:
			now := time.Now().Unix()
			world.Update(now)
			fps(e)
		}
	}
QUIT:
	e.Close()
}

// needs actual reload on restart, this loads once afaik
func reload(e *Engine, hefn HandleErrorFunc, llfn LoadLuaModuleFunc) {
	hefn(e, llfn())
}

var (
	frameDelta int64
	frameTime  int64
	frameCount int64
	FPS        int64
)

func fps(e *Engine) {
	now := time.Now().Unix()
	frameDelta = now - frameTime
	frameTime = now
	frameCount = frameCount + 1
	if frameDelta >= 1 {
		FPS = frameCount
		e.FPS = FPS
		frameCount = 0
	}
}

func (e *Engine) HandleError(r error) {
	e.hefn(e, r)
}

func (e *Engine) Close() {
	if e.last != nil {
		e.Print("closing with error")
		e.Print(e.last)
	} else {
		e.Print("closing....")
	}
	display.Close()
	e.Print("done")
	os.Exit(0)
}

func (e *Engine) Kill() {
	e.kill = true
}

func (e *Engine) KUEvent(k glfw.Key, m glfw.ModifierKey) {
	if e.debug {
		switch k {
		case glfw.KeyR:
			e.restart = true
		case glfw.KeyEscape:
			e.Kill()
		}
	}
}

type ConfigFn func(*Engine) error

type Config interface {
	Order() int
	Configure(*Engine) error
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

func (c config) Configure(e *Engine) error {
	return c.fn(e)
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
	e          *Engine
	configured bool
	list       configList
}

func newConfiguration(e *Engine, conf ...Config) *configuration {
	c := &configuration{
		e:    e,
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

func configure(e *Engine, conf ...Config) error {
	for _, c := range conf {
		err := c.Configure(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *configuration) Configure() error {
	sort.Sort(c.list)

	err := configure(c.e, c.list...)
	if err == nil {
		c.configured = true
	}

	return err
}

func (c *configuration) Configured() bool {
	return c.configured
}

var builtIns = []Config{
	config{1001, eLogger},
	config{1003, eError},
	config{5000, eGraphics},
	config{5001, eRenders},
	config{6000, eDisplay},
	config{7000, eInput},
	config{8001, eLua},
	config{8002, eCheckLoadLuaModule},
	config{9000, eWorld},
}

func eLogger(e *Engine) error {
	if e.Logger == nil {
		l := log.New(os.Stdout, log.LInfo, log.DefaultNullFormatter())
		log.Current = l
		e.Logger = l
	}
	return nil
}

func SetLogger(k string) Config {
	return NewConfig(2000,
		func(e *Engine) error {
			switch k {
			case "raw":
				e.SwapFormatter(log.GetFormatter("raw"))
			case "stdout", "text":
				e.SwapFormatter(log.GetFormatter("shiva_text"))
			}
			return nil
		})
}

func defaultHandleError(e *Engine, r error) {
	if r != nil {
		e.last = r
		e.kill = true
	}
}

func eError(e *Engine) error {
	if e.hefn == nil {
		e.hefn = defaultHandleError
	}
	return nil
}

func SetErrorHandler(fn HandleErrorFunc) Config {
	return NewConfig(50,
		func(e *Engine) error {
			e.hefn = fn
			return nil
		})
}

func eGraphics(e *Engine) error {
	if e.gp == nil {
		p, err := graphics.New("opengl4.5", e.debug)
		if err != nil {
			return err
		}
		e.gp = p
	}
	return nil
}

func SetGraphics(s string) Config {
	return NewConfig(50,
		func(e *Engine) error {
			gp, err := graphics.New(s, e.debug)
			if err != nil {
				return err
			}
			e.gp = gp
			return nil
		})
}

func eRenders(e *Engine) error {
	if e.rs == "" {
		e.rs = render.DefaultRenderer.String()
	}
	return nil
}

func SetRenders(s string) Config {
	return NewConfig(50,
		func(e *Engine) error {
			e.rs = s
			return nil
		})
}

func eDisplay(e *Engine) error {
	var err error
	if currentDisplaySystem == nil {
		currentDisplaySystem, err = display.New(e.gp, e.rs)
	}
	return err
}

//func SetDisplay(system) Config {
//	return NewConfig(50,
//		func(e *Engine) error {
//			return nil
//		})
//}

func eInput(e *Engine) error {
	w, err := display.Window()
	if err != nil {
		return err
	}
	input.Register(w)
	if e.debug {
		input.KeyUpInput.Subscribe(e)
	}

	currentInputSystem = input.CurrentInputSystem

	return nil
}

//func SetInput(system) Config {
//	return NewConfig(50,
//		func(e *Engine) error {
//			return nil
//		})
//}

var (
	luaDir  string = workingDir
	luaFile string = "main"
)

func SetLua(dir, file string) Config {
	return NewConfig(8000,
		func(e *Engine) error {
			luaDir = dir
			luaFile = file
			return nil
		})
}

var shv lua.Module = lua.ShvModule

func mkLoadLuaFn(L *lua.Lua, f string) func() error {
	return func() error {
		L.Push(L.RequireFn())
		L.PushString(f)
		if err := L.CallWithTraceback(1, 0); err != nil {
			return err
		}
		return nil
	}
}

func eLua(e *Engine) error {
	dlfn := display.RegisterWith()
	dlfn(shv)

	ilfn := input.RegisterWith()
	ilfn(shv)

	slfn := scene.RegisterWith()
	slfn(shv)

	L, err := lua.New(
		e.debug,
		lua.SetPath("_SHIVA_PATH", luaDir),
		lua.SetModules(math.Module),
	)
	if err != nil {
		return err
	}
	e.l = L
	e.llfn = mkLoadLuaFn(L, luaFile)
	return nil
}

var NoLoadLuaFuncError = xrror.Xrror("load lua function has not been specified.")

func eCheckLoadLuaModule(e *Engine) error {
	if e.llfn == nil {
		return NoLoadLuaFuncError
	}
	return nil
}

var (
	currentDisplaySystem ecs.System
	currentInputSystem   ecs.System
)

func eWorld(e *Engine) error {
	hefn := func(err error) {
		e.hefn(e, err)
	}
	world := ecs.New(hefn)
	world.Add(
		currentDisplaySystem,
		currentInputSystem,
	)
	e.w = world
	return nil
}

var workingDir string

func init() {
	workingDir, _ = os.Getwd()
}
