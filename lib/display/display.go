package display

import (
	"runtime"

	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/Laughs-In-Flowers/shiva/lib/math"
	"github.com/Laughs-In-Flowers/shiva/lib/render"
	"github.com/Laughs-In-Flowers/shiva/lib/scene"
	"github.com/Laughs-In-Flowers/shiva/lib/xrror"
	"github.com/go-gl/glfw/v3.2/glfw"

	l "github.com/yuin/gopher-lua"
)

var NoDisplay = xrror.Xrror("Display has not been created.")

func current() bool {
	return currentDisplay != nil
}

func Current() (*Display, error) {
	if current() {
		return currentDisplay, nil
	}
	return nil, NoDisplay
}

func Close() error {
	if current() {
		currentDisplay.SetShouldClose(true)
		currentDisplay.Destroy()
		glfw.Terminate()
		return nil
	}
	return NoDisplay
}

func Window() (*glfw.Window, error) {
	if current() {
		return currentDisplay.Window, nil
	}
	return nil, NoDisplay
}

var (
	currentDisplay   *Display
	currentDisplayId l.LValue
)

type DisplayData struct {
	*viewport
	title      string
	clearColor math.Vector
}

func DefaultDisplayData() *DisplayData {
	return &DisplayData{
		&viewport{0, 0, 640, 480},
		"shv",
		math.Vec4(0, 0, 0, 1),
	}
}

type viewport struct {
	x, y, w, h int
}

type Display struct {
	*glfw.Window
	*DisplayData
	*scene.Scene
}

var (
	defaultWidth  int = 640
	defaultHeight int = 480
)

func New(p graphics.Provider, renders string) (*Display, error) {
	var dis *Display
	var err error

	err = glfw.Init()
	if err != nil {
		return nil, err
	}

	err = p.Init()
	if err != nil {
		return nil, err
	}

	var nativeWindow *glfw.Window
	nativeWindow, err = glfw.CreateWindow(defaultWidth, defaultHeight, "", nil, nil)
	if err != nil {
		return nil, err
	}

	var r render.Renderer
	r, err = render.New(renders, p)
	if err != nil {
		return nil, err
	}

	_, major, minor := r.Version()
	glfw.WindowHint(glfw.ContextVersionMajor, major)
	glfw.WindowHint(glfw.ContextVersionMinor, minor)
	scene := scene.NewScene(r, nativeWindow)

	nativeWindow.MakeContextCurrent()
	glfw.SwapInterval(1)

	nativeWindow.SetPosCallback(dis.posCallback)
	nativeWindow.SetSizeCallback(dis.sizeCallback)
	nativeWindow.SetFramebufferSizeCallback(dis.framebufferSizeCallback)
	nativeWindow.SetCloseCallback(dis.closeCallback)
	nativeWindow.SetRefreshCallback(dis.refreshCallback)
	nativeWindow.SetFocusCallback(dis.focusCallback)
	nativeWindow.SetIconifyCallback(dis.iconifyCallback)

	dis = &Display{
		nativeWindow,
		DefaultDisplayData(),
		scene,
	}
	currentDisplay = dis
	return dis, err
}

func (d *Display) Priority() int {
	return 1
}

func (d *Display) Update(int64) error {
	d.Render()
	d.SwapBuffers()
	return nil
}

func (d *Display) Remove(uint64) {}

func (d *Display) posCallback(w *glfw.Window, xpos int, ypos int) {
	//spew.Dump("posCallback")
}

func (d *Display) sizeCallback(w *glfw.Window, width int, height int) {
	//notify scene in some way perhaps
	//spew.Dump("sizeCallback")
}

func (d *Display) framebufferSizeCallback(w *glfw.Window, width int, height int) {
	//spew.Dump("framebuffersizeCallback")
}

func (d *Display) closeCallback(w *glfw.Window) {
	//spew.Dump("closeCallback")
}

func (d *Display) refreshCallback(w *glfw.Window) {
	//spew.Dump("refreshCallback")
}

func (d *Display) focusCallback(w *glfw.Window, focused bool) {
	//spew.Dump("focusCallback")
}

func (d *Display) iconifyCallback(w *glfw.Window, iconified bool) {
	//spew.Dump("iconifyCallback")
}

func (d *Display) setTitle(t string) {
	d.title = t
	d.SetTitle(t)
}

func (d *Display) setWidth(wd int) {
	_, currentHeight := d.GetSize()
	d.SetSize(wd, currentHeight)
}

func (d *Display) setHeight(h int) {
	currentWidth, _ := d.GetSize()
	d.SetSize(currentWidth, h)
}

func (d *Display) setClearColor(v math.Vector) {
	d.ClearColor(v.Get(0), v.Get(1), v.Get(2), v.Get(3))
	d.clearColor = v
}

func (d *Display) pixelDensity() float64 {
	w := d.Window
	wwidth, _ := w.GetSize()
	fwidth, _ := w.GetFramebufferSize()
	return (float64(fwidth) / float64(wwidth))
}

const lDisplayClass = "WINDOW"

func lDisplay(L *l.LState) int {
	ud := L.NewUserData()
	ud.Value = currentDisplay
	switch L.GetTop() {
	case 1:
		setDisplayByTable(currentDisplay, L, 1)
	default:
		L.RaiseError("display takes only one table as an argument")
	}
	currentDisplayId = ud
	L.SetMetatable(ud, L.GetTypeMetatable(lDisplayClass))
	L.Push(ud)
	return 1
}

func setDisplayByTable(d *Display, L *l.LState, pos int) {
	t := L.CheckTable(pos)
	t.ForEach(func(k l.LValue, v l.LValue) {
		switch k.String() {
		case "width":
			//d.setWidth(l.LVAsNumber(v))
		case "height":
			//d.setHeight(l.LVAsNumber(v))
		case "title":
			d.setTitle(v.String())
		case "clear_color":
			d.setClearColor(math.ToVector(L, v))
		}
	})
}

func displayCall(t *lua.Table) l.LGFunction {
	return func(L *l.LState) int {
		setDisplayByTable(currentDisplay, L, 1)
		return 0
	}
}

func checkDisplay(L *l.LState) *Display {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Display); ok {
		return v
	}
	L.ArgError(1, "display expected")
	return nil
}

type displayMemberFunc func(*l.LState, *Display) int

func displayMember(fn displayMemberFunc) l.LGFunction {
	return func(L *l.LState) int {
		if display := checkDisplay(L); display != nil {
			return fn(L, display)
		}
		return 0
	}
}

func displayProperty(get, set displayMemberFunc) l.LGFunction {
	return lua.NewProperty(displayMember(get), displayMember(set))
}

func getTitle(L *l.LState, d *Display) int {
	L.Push(l.LString(d.title))
	return 1
}

func setTitle(L *l.LState, d *Display) int {
	n := L.Get(3)
	d.setTitle(n.String())
	return 0
}

func getWidth(L *l.LState, d *Display) int {
	wd, _ := d.GetSize()
	L.Push(l.LNumber(wd))
	return 1
}

func setWidth(L *l.LState, d *Display) int {
	width := L.CheckInt(2)
	d.setWidth(width)
	return 0
}

func getPixelWidth(L *l.LState, d *Display) int {
	pw, _ := d.GetFramebufferSize()
	L.Push(l.LNumber(pw))
	return 1
}

func getHeight(L *l.LState, d *Display) int {
	_, h := d.GetSize()
	L.Push(l.LNumber(h))
	return 1
}

func setHeight(L *l.LState, d *Display) int {
	height := L.CheckInt(2)
	d.setHeight(height)
	return 0
}

func getPixelHeight(L *l.LState, d *Display) int {
	_, ph := d.GetFramebufferSize()
	L.Push(l.LNumber(ph))
	return 1
}

func getLeft(L *l.LState, d *Display) int {
	left, _, _, _ := d.GetFrameSize()
	L.Push(l.LNumber(left))
	return 1
}

func getTop(L *l.LState, d *Display) int {
	_, top, _, _ := d.GetFrameSize()
	L.Push(l.LNumber(top))
	return 1
}

func getRight(L *l.LState, d *Display) int {
	_, _, right, _ := d.GetFrameSize()
	L.Push(l.LNumber(right))
	return 1
}

func getBottom(L *l.LState, d *Display) int {
	_, _, _, bottom := d.GetFrameSize()
	L.Push(l.LNumber(bottom))
	return 1
}

func getClearColor(L *l.LState, d *Display) int {
	fn := func(u *l.LUserData) {
		u.Value = d.clearColor
	}
	lua.PushNewUserData(L, fn, math.VEC4)
	return 1
}

func setClearColor(L *l.LState, d *Display) int {
	vec := math.UnpackToVec(L, 2, math.VEC4, true)
	d.setClearColor(vec)
	return 0
}

var displayTable = &lua.Table{
	lDisplayClass,
	nil,
	[]*lua.LMetaFunc{
		lua.DefaultIdx("__index"),
		lua.DefaultIdx("__newindex"),
	},
	map[string]l.LGFunction{
		"title":        displayProperty(getTitle, setTitle),
		"width":        displayProperty(getWidth, setWidth),
		"pixel_width":  displayProperty(getPixelWidth, nil),
		"height":       displayProperty(getHeight, setHeight),
		"pixel_height": displayProperty(getPixelHeight, nil),
		"left":         displayProperty(getLeft, nil),
		"top":          displayProperty(getTop, nil),
		"right":        displayProperty(getRight, nil),
		"bottom":       displayProperty(getBottom, nil),
		"clear_color":  displayProperty(getClearColor, setClearColor),
	},
	map[string]l.LGFunction{},
}

func RegisterWith() lua.RegisterWith {
	return func(m lua.Module) error {
		m.AddLGFunc("display", lDisplay)
		rmtfn := func(L *l.LState, M lua.Module) {
			M.Register(L, displayTable)
		}
		m.AddMT(rmtfn)
		return nil
	}
}

func init() {
	runtime.LockOSThread()
}
