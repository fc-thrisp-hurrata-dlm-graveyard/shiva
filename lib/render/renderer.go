package render

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/shader"
	"github.com/Laughs-In-Flowers/shiva/lib/math"
	"github.com/Laughs-In-Flowers/shiva/lib/xrror"
)

type RendererT int

func (t RendererT) String() string {
	switch t {
	case FORWARD:
		return "forward"
	}
	return "UNKNOWN"
}

func StringToRendererT(s string) RendererT {
	switch s {
	case "forward":
		return FORWARD
	}
	return UNKNOWN
}

const (
	UNKNOWN RendererT = iota
	FORWARD
)

var DefaultRenderer = FORWARD

var UnknownRendererError = xrror.Xrror("%s is not a known renderer").Out

type NewRendererFunc func(graphics.Provider) Renderer

type rpr struct {
	has map[RendererT]NewRendererFunc
}

func (r *rpr) add(name string, fn NewRendererFunc) {
	rn := StringToRendererT(name)
	r.has[rn] = fn
}

func (r *rpr) get(name string, p graphics.Provider) (Renderer, error) {
	var rr Renderer
	rt := StringToRendererT(name)
	if rfn, exists := r.has[rt]; exists {
		rr = rfn(p)
		return rr, nil
	}
	return nil, UnknownRendererError(name)
}

var RendererRegistry *rpr

func Register(name string, fn NewRendererFunc) {
	RendererRegistry.add(name, fn)
}

func New(s string, p graphics.Provider) (Renderer, error) {
	return RendererRegistry.get(s, p)
}

type Renderable interface {
	Renderable() bool
	SetRenderable(bool)
	Render(Renderer)
}

type Space interface {
	ViewMatrice() math.Matrice
	SetViewMatrice(math.Matrice)
	ProjectionMatrice() math.Matrice
	SetProjectionMatrice(math.Matrice)
	Last() math.Matrice
	SetLast(math.Matrice)
}

type Renderer interface {
	graphics.Provider
	shader.Shaderer
	Space
	Type() RendererT
	Initialize()
	Rend(...Renderable)
}

func init() {
	RendererRegistry = &rpr{make(map[RendererT]NewRendererFunc)}
	Register("forward", newForwardRenderer)
}
