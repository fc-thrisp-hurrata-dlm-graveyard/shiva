package render

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/shader"
)

//type Params struct {
//vec map[string]Vector
//mat map[string]Matrice
//flo map[string]float32
//}

//func NewParams() *Params {
//return &Params{
//make(map[string]Vector),
//make(map[string]Matrice),
//make(map[string]float32),
//}
//}

//func (p *Params) Vector(k string) Vector {
//	return nil
//}

//func (p *Params) Matrice(k string) Matrice {
//	return nil
//}

//func (p *Params) Float(k string) float32 {
//	return 0
//}

type fRenderer struct {
	graphics.Provider
	shader.Shaderer
	//pass, nextPass, passMask uint32
	//boundViewport            *viewportState
	//activeViewport           *viewportState
	//boundColorMask           *colorMaskState
	//activeColorMask          *colorMaskState
	//boundDepth               *depthState
	//activeDepth              *depthState
	//boundStencil             *stencilState
	//activeStencil            *stencilState
	//boundScissor             *scissorState
	//activeScissor            *scissorState
	//boundCullface            *cullFaceState
	//activeCullFace           *cullFaceState
	//boundBlend               *blendState
	//activeBlend              *blendState
	//maxDrawArraySize         int
	//numberEnabledVaas        int
	//params                   *Params
}

func newForwardRenderer(gp graphics.Provider) Renderer {
	return &fRenderer{
		gp,
		shader.DefaultShaderer(),
		//1, 1, 1,
		//nil, nil,
		//nil, nil,
		//nil, nil,
		//nil, nil,
		//nil, nil,
		//nil, nil,
		//nil, nil,
		//0,
		//0,
		//NewParams(),
	}
}

func (r *fRenderer) Type() RendererT {
	return FORWARD
}

func (r *fRenderer) Rend(d ...Renderable) {
	r.pre()
	for _, n := range d {
		n.Render(r)
	}
	r.post()
}

func (r *fRenderer) pre() {
	r.Clear(graphics.COLOR_BUFFER_BIT | graphics.DEPTH_BUFFER_BIT | graphics.STENCIL_BUFFER_BIT)
}

func (r *fRenderer) post() {
	r.UseProgram(0) // unbind current shader program
}

//boundProgramId
//activeProgram
//DrawArrays()
//DrawElements()
//ValidateActiveProgram()
//BindActiveProgram
//BindActiveProgramParams
//UpdateState
//EnableVaas
//DoRender

//viewport state
//type viewportState struct {
//	x, y, w, h int
//}

//func (v *viewportState) set(x, y, w, h int) {
//	v.x = x
//	v.y = y
//	v.w = w
//	v.h = h
//}

//func (v *viewportState) restore(old *viewportState) {
//	v.set(old.x, old.y, old.w, old.h)
//}

//func (v *viewportState) bind(*Renderer) {}

//color mask state
//type colorMaskState struct {
//	r, g, b, a bool
//}

//func (c *colorMaskState) set() {}

//func (c *colorMaskState) restore(old *colorMaskState) {}

//func (c *colorMaskState) bind(*Renderer) {}

//depth test state
//type depthState struct {
//testEnabled bool
//maskEnabled bool
//fn
//}

//blend state
//type blendState struct {
//enabled bool
//equation_rgb    = AM_BLEND_EQUATION_ADD;
//equation_alpha  = AM_BLEND_EQUATION_ADD;
//sfactor_rgb     = AM_BLEND_SFACTOR_SRC_ALPHA;
//dfactor_rgb     = AM_BLEND_DFACTOR_ONE_MINUS_SRC_ALPHA;
//sfactor_alpha   = AM_BLEND_SFACTOR_SRC_ALPHA;
//dfactor_alpha   = AM_BLEND_DFACTOR_ONE_MINUS_SRC_ALPHA;
//constantR, constantG, constantB, constantA float32
//}

//stencil test state
//type stencilState struct{}

//scissor test state
//type scissorState struct {
//	enabled    bool
//	x, y, w, h int
//}

//func (s *scissorState) set() {}

//func (s *scissorState) restore(old *scissorState) {}

//func (s *scissorState) bind(*Renderer) {}

//cull face state
//type cullFaceState struct{}

//sample coverage state

//polygon_offset_state

//line_state

//dither_state

//func (r *Renderer) ErrCount() int {
//	return len(r.errs)
//}

//func (r *Renderer) SetMaxErrCount(c int) {
//	r.errOn = c
//}

//var TooManyErrors = Xrror("Too many render errors (%d), last: %s").Out

//func (r *Renderer) Render(n Node) error {
//	s := &State{}
//	err := n.Render(s)
//	if err != nil {
//		r.errs = append(r.errs, err)
//	}
//	if r.errOn != 0 && len(r.errs) >= r.errOn {
//		last := r.errs[len(r.errs)-1]
//		return TooManyErrors(r.errOn, last)
//	}
//	s = nil
//	return nil
//}
