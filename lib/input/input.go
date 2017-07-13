package input

import (
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func Register(w *glfw.Window) {
	glfw.SetJoystickCallback(JoystickConnectInput.Callback())
	w.SetKeyCallback(KeyInput.Callback())
	w.SetCharCallback(CharInput.Callback())
	w.SetCharModsCallback(CharModInput.Callback())
	w.SetMouseButtonCallback(MouseButtonInput.Callback())
	w.SetCursorPosCallback(CursorPositionInput.Callback())
	w.SetCursorEnterCallback(CursorEnterInput.Callback())
	w.SetScrollCallback(ScrollInput.Callback())
	w.SetDropCallback(DropInput.Callback())
}

var JoystickConnectInput *joystickConnectInput

type jcsubscriber interface {
	JCEvent(int, int)
}

type joystickConnectInput struct {
	subscribers []jcsubscriber
}

func (j *joystickConnectInput) Callback() glfw.JoystickCallback {
	return func(joy, event int) {
		for _, s := range j.subscribers {
			s.JCEvent(joy, event)
		}
	}
}

func (j *joystickConnectInput) Subscribe(ss ...jcsubscriber) {
	if j.subscribers == nil {
		j.subscribers = make([]jcsubscriber, 0)
	}
	j.subscribers = append(j.subscribers, ss...)
}

// joystick axis / button make on plugin --> subscribe

var KeyInput *keyInput

type ksubscriber interface {
	KEvent(glfw.Key, int, glfw.Action, glfw.ModifierKey)
}

type keyInput struct {
	subscribers []ksubscriber
}

func (ki *keyInput) Callback() glfw.KeyCallback {
	return func(
		w *glfw.Window,
		k glfw.Key,
		s int,
		a glfw.Action,
		m glfw.ModifierKey,
	) {
		for _, sb := range ki.subscribers {
			sb.KEvent(k, s, a, m)
		}
	}
}

func (ki *keyInput) Subscribe(ss ...ksubscriber) {
	if ki.subscribers == nil {
		ki.subscribers = make([]ksubscriber, 0)
	}
	ki.subscribers = append(ki.subscribers, ss...)
}

var KeyDownInput *keyDownInput

type keyDownSubscriber interface {
	KDEvent(glfw.Key, glfw.ModifierKey)
}

type keyDownInput struct {
	subscribers []keyDownSubscriber
}

func (kd *keyDownInput) KEvent(k glfw.Key, s int, a glfw.Action, m glfw.ModifierKey) {
	if a == glfw.Press {
		for _, sb := range kd.subscribers {
			sb.KDEvent(k, m)
		}
	}
}

func (kd *keyDownInput) Subscribe(ss ...keyDownSubscriber) {
	if kd.subscribers == nil {
		kd.subscribers = make([]keyDownSubscriber, 0)
	}
	kd.subscribers = append(kd.subscribers, ss...)
}

var KeyUpInput *keyUpInput

type keyUpSubscriber interface {
	KUEvent(glfw.Key, glfw.ModifierKey)
}

type keyUpInput struct {
	subscribers []keyUpSubscriber
}

func (ku *keyUpInput) KEvent(k glfw.Key, s int, a glfw.Action, m glfw.ModifierKey) {
	if a == glfw.Release {
		for _, sb := range ku.subscribers {
			sb.KUEvent(k, m)
		}
	}
}

func (ku *keyUpInput) Subscribe(ss ...keyUpSubscriber) {
	if ku.subscribers == nil {
		ku.subscribers = make([]keyUpSubscriber, 0)
	}
	ku.subscribers = append(ku.subscribers, ss...)
}

var CharInput *charInput

type cisubscriber interface {
	CIEvent(rune)
}

type charInput struct {
	subscribers []cisubscriber
}

func (c *charInput) Callback() glfw.CharCallback {
	return func(w *glfw.Window, r rune) {
		for _, s := range c.subscribers {
			s.CIEvent(r)
		}
	}
}

func (c *charInput) Subscribe(ss ...cisubscriber) {
	if c.subscribers == nil {
		c.subscribers = make([]cisubscriber, 0)
	}
	c.subscribers = append(c.subscribers, ss...)
}

var CharModInput *charModInput

type cmsubscriber interface {
	CMEvent(rune, glfw.ModifierKey)
}

type charModInput struct {
	subscribers []cmsubscriber
}

func (c *charModInput) Callback() glfw.CharModsCallback {
	return func(w *glfw.Window, r rune, m glfw.ModifierKey) {
		for _, s := range c.subscribers {
			s.CMEvent(r, m)
		}
	}
}

func (c *charModInput) Subscribe(ss ...cmsubscriber) {
	if c.subscribers == nil {
		c.subscribers = make([]cmsubscriber, 0)
	}
	c.subscribers = append(c.subscribers, ss...)
}

var MouseButtonInput *mouseButtonInput

type mbsubscriber interface {
	MBEvent(glfw.MouseButton, glfw.Action, glfw.ModifierKey)
}

type mouseButtonInput struct {
	subscribers []mbsubscriber
}

func (mb *mouseButtonInput) Callback() glfw.MouseButtonCallback {
	return func(
		w *glfw.Window,
		b glfw.MouseButton,
		a glfw.Action,
		m glfw.ModifierKey,
	) {
		for _, s := range mb.subscribers {
			s.MBEvent(b, a, m)
		}
	}
}

func (mb *mouseButtonInput) Subscribe(ss ...mbsubscriber) {
	if mb.subscribers == nil {
		mb.subscribers = make([]mbsubscriber, 0)
	}
	mb.subscribers = append(mb.subscribers, ss...)
}

var CursorPositionInput *cursorPositionInput

type cpsubscriber interface {
	CPEvent(float64, float64)
}

type cursorPositionInput struct {
	subscribers []cpsubscriber
}

func (c *cursorPositionInput) Callback() glfw.CursorPosCallback {
	return func(w *glfw.Window, xpos float64, ypos float64) {
		for _, s := range c.subscribers {
			s.CPEvent(xpos, ypos)
		}
	}
}

func (c *cursorPositionInput) Subscribe(ss ...cpsubscriber) {
	if c.subscribers == nil {
		c.subscribers = make([]cpsubscriber, 0)
	}
	c.subscribers = append(c.subscribers, ss...)
}

var CursorEnterInput *cursorEnterInput

type cesubscriber interface {
	CEEvent(bool)
}

type cursorEnterInput struct {
	subscribers []cesubscriber
}

func (c *cursorEnterInput) Callback() glfw.CursorEnterCallback {
	return func(w *glfw.Window, entered bool) {
		for _, s := range c.subscribers {
			s.CEEvent(entered)
		}
	}
}

func (c *cursorEnterInput) Subscribe(ss ...cesubscriber) {
	if c.subscribers == nil {
		c.subscribers = make([]cesubscriber, 0)
	}
	c.subscribers = append(c.subscribers, ss...)
}

var ScrollInput *scrollInput

type ssubscriber interface {
	SEvent(float64, float64)
}

type scrollInput struct {
	subscribers []ssubscriber
}

func (s *scrollInput) Callback() glfw.ScrollCallback {
	return func(w *glfw.Window, xoff float64, yoff float64) {
		for _, sb := range s.subscribers {
			sb.SEvent(xoff, yoff)
		}
	}
}

func (s *scrollInput) Subscribe(ss ...ssubscriber) {
	if s.subscribers == nil {
		s.subscribers = make([]ssubscriber, 0)
	}
	s.subscribers = append(s.subscribers, ss...)
}

var DropInput *dropInput

type dsubscriber interface {
	DEvent([]string)
}

type dropInput struct {
	subscribers []dsubscriber
}

func (d *dropInput) Callback() glfw.DropCallback {
	return func(w *glfw.Window, paths []string) {
		for _, s := range d.subscribers {
			s.DEvent(paths)
		}
	}
}

func (d *dropInput) Subscribe(ss ...dsubscriber) {
	if d.subscribers == nil {
		d.subscribers = make([]dsubscriber, 0)
	}
	d.subscribers = append(d.subscribers, ss...)
}

// input event poll
// secondary input systems, e.g. input subscribers fan out, queuing, reporting
var CurrentInputSystem *inputSystem

type inputSystem struct{}

func (is *inputSystem) Priority() int {
	return 2
}

func (is *inputSystem) Update(d int64) error {
	glfw.PollEvents()
	return nil
}

func (is *inputSystem) Remove(uint64) {}

func RegisterWith() lua.RegisterWith {
	return func(m lua.Module) error {
		return nil
	}
}

func init() {
	JoystickConnectInput = &joystickConnectInput{}
	KeyInput = &keyInput{}
	KeyDownInput = &keyDownInput{}
	KeyUpInput = &keyUpInput{}
	KeyInput.Subscribe(KeyDownInput, KeyUpInput)
	CharInput = &charInput{}
	CharModInput = &charModInput{}
	MouseButtonInput = &mouseButtonInput{}
	CursorPositionInput = &cursorPositionInput{}
	CursorEnterInput = &cursorEnterInput{}
	ScrollInput = &scrollInput{}
	DropInput = &dropInput{}

	CurrentInputSystem = &inputSystem{}
}
