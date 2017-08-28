package render

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/shader"
	"github.com/Laughs-In-Flowers/shiva/lib/math"
)

type fRenderer struct {
	graphics.Provider
	shader.Shaderer
	view math.Matrice
	proj math.Matrice
	last math.Matrice
}

func newForwardRenderer(gp graphics.Provider) Renderer {
	r := &fRenderer{
		gp,
		shader.DefaultShaderer(),
		math.Mat4(), math.Mat4(), math.Mat4(),
	}
	r.Initialize()
	return r
}

func (r *fRenderer) ViewMatrice() math.Matrice {
	return r.view
}

func (r *fRenderer) SetViewMatrice(m math.Matrice) {
	r.view = m
}

func (r *fRenderer) ProjectionMatrice() math.Matrice {
	return r.proj
}

func (r *fRenderer) SetProjectionMatrice(m math.Matrice) {
	r.proj = m
}

func (r *fRenderer) Last() math.Matrice {
	return r.last
}

func (r *fRenderer) SetLast(m math.Matrice) {
	r.last = m
}

func (r *fRenderer) Type() RendererT {
	return FORWARD
}

func (r *fRenderer) Initialize() {
	r.Clear(graphics.COLOR_BUFFER_BIT | graphics.DEPTH_BUFFER_BIT | graphics.STENCIL_BUFFER_BIT)
	r.ClearDepth(1)
	r.ClearStencil(0)
	r.Enable(graphics.DEPTH_TEST)
	r.DepthFunc(graphics.LEQUAL)
	r.FrontFace(graphics.CCW)
	r.CullFace(graphics.BACK)
	r.Enable(graphics.CULL_FACE)
	r.Enable(graphics.BLEND)
	r.BlendEquation(graphics.FUNC_ADD)
	r.BlendFunc(graphics.SRC_ALPHA, graphics.ONE_MINUS_SRC_ALPHA)
	r.Enable(graphics.VERTEX_PROGRAM_POINT_SIZE)
	r.Enable(graphics.PROGRAM_POINT_SIZE)
	r.Enable(graphics.MULTISAMPLE)
	r.Enable(graphics.POLYGON_OFFSET_FILL)
	r.Enable(graphics.POLYGON_OFFSET_LINE)
	r.Enable(graphics.POLYGON_OFFSET_POINT)
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
	r.UseProgram(0)
}

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
