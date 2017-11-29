package scene

import (
	"sync"

	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/Laughs-In-Flowers/shiva/lib/render"
	"github.com/go-gl/glfw/v3.2/glfw"

	l "github.com/yuin/gopher-lua"
)

var currentScene *Scene

type nenderable struct {
	sync.RWMutex
	nodes []Node
}

func newNenderable() *nenderable {
	return &nenderable{
		sync.RWMutex{},
		make([]Node, 0),
	}
}

func (n *nenderable) Count() int {
	n.Lock()
	count := 0
	for _, n := range n.nodes {
		count = count + n.Count()
	}
	n.Unlock()
	return count
}

func (n *nenderable) Renderable() bool {
	return true
}

func (n *nenderable) SetRenderable(bool) {}

func (n *nenderable) Render(r render.Renderer) {
	for _, n := range n.nodes {
		n.Render(r)
	}
}

func (n *nenderable) attach(d Node) {
	n.RLock()
	n.nodes = append(n.nodes, d)
	n.RUnlock()
}

func (n *nenderable) detach(d Node) {
	for idx, nn := range n.nodes {
		if nn.Equal(d) {
			n.RLock()
			n.nodes = n.nodes[:idx+copy(n.nodes[idx:], n.nodes[idx+1:])]
			n.RUnlock()
		}
	}
}

func (n *nenderable) clear() {
	n.RLock()
	n.nodes = make([]Node, 0)
	n.RUnlock()
}

func (n *nenderable) List() []Node {
	return n.nodes
}

type Scene struct {
	render.Renderer
	n      *nenderable
	update bool
}

var nativeWindow *glfw.Window

func Aspect(w *glfw.Window) float32 {
	var ret float32
	if w != nil {
		width, height := w.GetSize()
		ret = float32(width) / float32(height)
	}
	return ret
}

func NewScene(r render.Renderer, nw *glfw.Window) *Scene {
	s := &Scene{
		r,
		newNenderable(),
		true,
	}
	nativeWindow = nw
	currentScene = s
	return s
}

func (s *Scene) Render() {
	r := s.Renderer
	r.Rend(s.n)
}

func (s *Scene) Attach(ns ...Node) {
	for _, n := range ns {
		s.n.attach(n)
	}
	s.update = true
}

func (s *Scene) Detach(ns ...Node) {
	for _, n := range ns {
		s.n.detach(n)
	}
	s.update = true
}

func (s *Scene) Clear() {
	s.n.clear()
}

func (s *Scene) Count() int {
	return s.n.Count()
}

func (s *Scene) Updateable() bool {
	return s.update
}

const lSceneClass = "SCENE"

func lScene(L *l.LState) int {
	ud := L.NewUserData()
	ud.Value = currentScene
	L.SetMetatable(ud, L.GetTypeMetatable(lSceneClass))
	L.Push(ud)
	return 1
}

func checkScene(L *l.LState, pos int) (*l.LUserData, *Scene) {
	ud := L.CheckUserData(pos)
	if n, ok := ud.Value.(*Scene); ok {
		return ud, n
	}
	L.ArgError(pos, "scene expected")
	return nil, nil
}

type sceneMemberFunc func(*l.LState, *l.LUserData, *Scene) int

func sceneMember(fn sceneMemberFunc) l.LGFunction {
	return func(L *l.LState) int {
		if u, n := checkScene(L, 1); n != nil {
			return fn(L, u, n)
		}
		return 0
	}
}

func sceneProperty(get, set sceneMemberFunc) l.LGFunction {
	return lua.NewProperty(sceneMember(get), sceneMember(set))
}

func getNodeCount(L *l.LState, u *l.LUserData, s *Scene) int {
	L.Push(l.LNumber(s.Count()))
	return 1
}

func attachNode(L *l.LState, u *l.LUserData, s *Scene) int {
	if n := pullNode(L, 2); n != nil {
		s.Attach(n)
		L.Push(l.LBool(true))
		return 1
	}
	return 0
}

func detachNode(L *l.LState, u *l.LUserData, s *Scene) int {
	if ud, n := pullNodeUD(L, 2); n != nil {
		s.Detach(n)
		L.Push(ud)
		return 1
	}
	return 0
}

func clearScene(L *l.LState, u *l.LUserData, s *Scene) int {
	s.Clear()
	return 0
}

var sceneTable = &lua.Table{
	lSceneClass,
	nil,
	[]*lua.LMetaFunc{
		lua.DefaultIdx("__index"),
		lua.DefaultIdx("__newindex"),
	},
	map[string]l.LGFunction{
		"count": sceneProperty(getNodeCount, nil),
	},
	map[string]l.LGFunction{
		"attach": sceneMember(attachNode),
		"detach": sceneMember(detachNode),
		"clear":  sceneMember(clearScene),
	},
}

type sceneRegister struct {
	has []sceneRegisterFunc
}

func (s *sceneRegister) add(fn sceneRegisterFunc) {
	s.has = append(s.has, fn)
}

func (s *sceneRegister) run(m lua.Module) error {
	var err error
	for _, fn := range s.has {
		err = fn(m)
		if err != nil {
			return err
		}
	}
	return err
}

type sceneRegisterFunc func(lua.Module) error

func registerWith(
	tag string,
	rfn l.LGFunction,
	tbl *lua.Table,
) sceneRegisterFunc {
	return func(mod lua.Module) error {
		if tag != "" && rfn != nil {
			mod.AddLGFunc(tag, rfn)
		}
		if tbl != nil {
			rmtfn := func(L *l.LState, M lua.Module) {
				M.Register(L, tbl)
			}
			mod.AddMT(rmtfn)
		}
		return nil
	}
}

func RegisterWith() lua.RegisterWith {
	return func(m lua.Module) error {
		sr := &sceneRegister{make([]sceneRegisterFunc, 0)}
		sr.add(registerWith("scene", lScene, sceneTable))
		//sr.add(registerWith("", nil, nodeTable))
		sr.add(registerWith("dummy", ldummy, defaultNodeTable(lDummyNodeClass)))
		sr.add(registerWith("root", lroot, defaultNodeTable(lRootNodeClass)))
		sr.add(registerWith("group", lgroup, defaultNodeTable(
			lGroupNodeClass,
			&lua.LMetaFunc{"__call", groupNodeCall},
		)))
		sr.add(registerWith("translate", ltranslate, lTranslateNodeTable))
		sr.add(registerWith("scale", lscale, lScaleNodeTable))
		sr.add(registerWith("rotate", lrotate, lRotateNodeTable))
		sr.add(registerWith("axis", laxis, lAxisNodeTable))
		sr.add(registerWith("sphere", lsphere, lSphereNodeTable))
		sr.add(registerWith("camera", lcamera, lCameraNodeTable))
		// default orthographic camera
		// default perspective camera
		return sr.run(m)
	}
}
