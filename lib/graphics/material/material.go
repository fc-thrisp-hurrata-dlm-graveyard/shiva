package material

import (
	"reflect"

	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/texture"
)

type Material interface {
	Materializer
	Lighter
	Sider
	Blender
	Depther
	Liner
	Framer
	Polygoner
	Texturer
}

type UseLights int

const (
	ULNone        UseLights = 0x00
	ULAmbient     UseLights = 0x01
	ULDirectional UseLights = 0x02
	ULPoint       UseLights = 0x04
	ULSpot        UseLights = 0x08
	ULAll         UseLights = 0xFF
)

type Side int

const (
	SIFront Side = iota
	SIBack
	SIDouble
)

type Blending int

const (
	BLNone Blending = iota
	BLNormal
	BLAdditive
	BLSubtractive
	BLMultiply
	BLCustom
)

type material struct {
	refcount         int               // Current number of references
	useShader        string            // Shader name
	independent      bool              // shader does not depend on the number of lights in the scene and/or number of textures in the material.
	uselights        UseLights         // Use lights bit mask
	sideVisible      Side              // sides visible
	wireframe        bool              // show as wirefrme
	depthMask        bool              // Enable writing into the depth buffer
	depthTest        bool              // Enable depth buffer test
	depthFunc        graphics.Enum     // Active depth test function
	blending         Blending          // blending mode
	blendRGB         graphics.Enum     // separate blend equation for RGB
	blendAlpha       graphics.Enum     // separate blend equation for Alpha
	blendSrcRGB      graphics.Enum     // separate blend func source RGB
	blendDstRGB      graphics.Enum     // separate blend func dest RGB
	blendSrcAlpha    graphics.Enum     // separate blend func source Alpha
	blendDstAlpha    graphics.Enum     // separate blend func dest Alpha
	lineWidth        float32           // line width for lines and mesh wireframe
	polyOffsetFactor float32           // polygon offset factor
	polyOffsetUnits  float32           // polygon offset units
	textures         []texture.Texture // List of textures
}

func New() *material {
	m := &material{}
	m.Initialize()
	return m
}

func Basic() Material {
	b := New()
	b.SetShader("basic")
	b.SetIndependent(true)
	return b
}

type Materializer interface {
	graphics.Initializer
	graphics.Closer
	graphics.Providable
	graphics.RefCounter
	//graphics.Renderable
	Shaderer
}

func (m *material) Initialize() {
	m.refcount = 1
	m.uselights = ULAll
	m.sideVisible = SIFront
	m.wireframe = false
	m.depthMask = true
	m.depthFunc = graphics.LEQUAL
	m.depthTest = true
	m.blending = BLNormal
	m.lineWidth = 1.0
	m.polyOffsetFactor = 0
	m.polyOffsetUnits = 0
	m.textures = make([]texture.Texture, 0)
}

func (m *material) Close() {
	if m.refcount > 1 {
		m.refcount--
		return
	}
	for i := 0; i < len(m.textures); i++ {
		m.textures[i].Close()
	}
	m.Initialize()
}

func (m *material) Providable() bool {
	return true
}

func (m *material) SetProvidable(b bool) {}

func (m *material) Provide(p graphics.Provider) {
	switch m.sideVisible {
	case SIFront:
		p.Enable(graphics.CULL_FACE)
		p.FrontFace(graphics.CCW)
	case SIBack:
		p.Enable(graphics.CULL_FACE)
		p.FrontFace(graphics.CW)
	case SIDouble:
		p.Disable(graphics.CULL_FACE)
		p.FrontFace(graphics.CCW)
	}

	if m.depthTest {
		p.Enable(graphics.DEPTH_TEST)
	} else {
		p.Disable(graphics.DEPTH_TEST)
	}
	p.DepthMask(m.depthMask)
	p.DepthFunc(m.depthFunc)

	if m.wireframe {
		p.PolygonMode(graphics.FRONT_AND_BACK, graphics.LINE)
	} else {
		p.PolygonMode(graphics.FRONT_AND_BACK, graphics.FILL)
	}

	// Set polygon offset if requested
	p.PolygonOffset(m.polyOffsetFactor, m.polyOffsetUnits)

	// Sets line width
	p.LineWidth(m.lineWidth)

	// Sets blending
	switch m.blending {
	case BLNone:
		p.Disable(graphics.BLEND)
	case BLNormal:
		p.Enable(graphics.BLEND)
		p.BlendEquationSeparate(graphics.FUNC_ADD, graphics.FUNC_ADD)
		p.BlendFunc(graphics.SRC_ALPHA, graphics.ONE_MINUS_SRC_ALPHA)
	case BLAdditive:
		p.Enable(graphics.BLEND)
		p.BlendEquation(graphics.FUNC_ADD)
		p.BlendFunc(graphics.SRC_ALPHA, graphics.ONE)
	case BLSubtractive:
		p.Enable(graphics.BLEND)
		p.BlendEquation(graphics.FUNC_ADD)
		p.BlendFunc(graphics.ZERO, graphics.ONE_MINUS_SRC_COLOR)
	case BLMultiply:
		p.Enable(graphics.BLEND)
		p.BlendEquation(graphics.FUNC_ADD)
		p.BlendFunc(graphics.ZERO, graphics.SRC_COLOR)
	case BLCustom:
		p.BlendEquationSeparate(m.blendRGB, m.blendAlpha)
		p.BlendFuncSeparate(m.blendSrcRGB, m.blendDstRGB, m.blendSrcAlpha, m.blendDstAlpha)
	default:
		panic("Invalid blending")
	}

	for idx, tex := range m.textures {
		tex.Render(p, idx)
	}
}

func (m *material) Increment() {
	m.refcount++
}

func (m *material) Decrement() {
	m.refcount--
}

type Shaderer interface {
	Shader() string
	SetShader(string)
	Independent() bool
	SetIndependent(bool)
}

func (m *material) Shader() string {
	return m.useShader
}

func (m *material) SetShader(s string) {
	m.useShader = s
}

func (m *material) Independent() bool {
	return m.independent
}

func (m *material) SetIndependent(b bool) {
	m.independent = b
}

type Lighter interface {
	UseLights() UseLights
	SetUseLights(UseLights)
}

func (m *material) UseLights() UseLights {
	return m.uselights
}

func (m *material) SetUseLights(u UseLights) {
	m.uselights = u
}

type Sider interface {
	Side() Side
	SetSide(Side)
}

func (m *material) Side() Side {
	return m.sideVisible
}

func (m *material) SetSide(s Side) {
	m.sideVisible = s
}

type Blender interface {
	SetBlending(Blending)
}

func (m *material) SetBlending(b Blending) {
	m.blending = b
}

type Depther interface {
	SetDepthMask(bool)
	SetDepthTest(bool)
}

func (m *material) SetDepthMask(state bool) {
	m.depthMask = state
}

func (m *material) SetDepthTest(state bool) {
	m.depthTest = state
}

type Liner interface {
	SetLineWidth(float32)
}

func (m *material) SetLineWidth(w float32) {
	m.lineWidth = w
}

type Framer interface {
	SetWireframe(bool)
}

func (m *material) SetWireframe(state bool) {
	m.wireframe = state
}

type Polygoner interface {
	SetPolygonOffset(float32, float32)
}

func (m *material) SetPolygonOffset(factor, units float32) {
	m.polyOffsetFactor = factor
	m.polyOffsetUnits = units
}

type Texturer interface {
	AddTexture(...texture.Texture)
	RemoveTexture(...texture.Texture)
	HasTexture(texture.Texture) bool
	TextureCount() int
}

func (m *material) AddTexture(t ...texture.Texture) {
	for _, tx := range t {
		m.textures = append(m.textures, tx)
	}
}

// confirm this works...
func Equals(t1, t2 texture.Texture) bool {
	return reflect.DeepEqual(t1, t2)
}

func (m *material) removeTexture(tx texture.Texture) {
	for idx, tex := range m.textures {
		if Equals(tex, tx) {
			copy(m.textures[idx:], m.textures[idx+1:])
			m.textures[len(m.textures)-1] = nil
			m.textures = m.textures[:len(m.textures)-1]
			break
		}
	}
}

func (m *material) RemoveTexture(t ...texture.Texture) {
	for _, tx := range t {
		m.removeTexture(tx)
	}
}

func (m *material) HasTexture(t texture.Texture) bool {
	for _, tx := range m.textures {
		if Equals(tx, t) {
			return true
		}
	}
	return false
}

func (m *material) TextureCount() int {
	return len(m.textures)
}
