package render

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/shader"
	"github.com/Laughs-In-Flowers/shiva/lib/xrror"
)

type RendererT int

func (t RendererT) String() string {
	switch t {
	case FORWARD:
		return "forward"
	case DEFERRED:
		return "deferred"
	}
	return "UNKNOWN"
}

func StringToRendererT(s string) RendererT {
	switch s {
	case "forward":
		return FORWARD
	case "deferred":
		return DEFERRED
	}
	return UNKNOWN
}

const (
	FORWARD RendererT = iota
	DEFERRED
	UNKNOWN
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
	Render(Renderer)
}

type Renderer interface {
	graphics.Provider
	shader.Shaderer
	Type() RendererT
	Rend(...Renderable)
}

func init() {
	RendererRegistry = &rpr{make(map[RendererT]NewRendererFunc)}
	Register("forward", newForwardRenderer)
	//Register("deferred", newDeferredRenderer)
}
