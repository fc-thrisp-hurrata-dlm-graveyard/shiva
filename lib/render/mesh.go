package render

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/geometry"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/material"
	"github.com/Laughs-In-Flowers/shiva/lib/xrror"
)

type Mesh interface {
	graphics.Initializer
	graphics.Closer
	graphics.Moder
	Geometer
	Materializer
	Renderable
}

type innerRenderFunc func(Renderer)

type mesh struct {
	tag        string
	g          geometry.Geometry
	materials  []Material
	mode       graphics.Enum
	renderable bool
	rfn        innerRenderFunc
}

func NewMesh(tag string, e geometry.Geometry, rfn innerRenderFunc, mode graphics.Enum) *mesh {
	m := &mesh{
		tag: tag,
		rfn: rfn,
	}
	m.Initialize()
	m.g = e
	m.mode = mode
	return m
}

func (m *mesh) Initialize() {
	m.g = nil
	m.materials = make([]Material, 0)
	m.renderable = true
}

func (m *mesh) Close() {
	m.g.Close()
	for i := 0; i < len(m.materials); i++ {
		a := m.materials[i].m
		a.Close()
	}
}

func (m *mesh) Mode() graphics.Enum {
	return m.mode
}

func (m *mesh) SetMode(e graphics.Enum) {
	m.mode = e
}

func (m *mesh) Renderable() bool {
	return m.renderable
}

func (m *mesh) SetRenderable(as bool) {
	m.renderable = as
}

func (m *mesh) Render(r Renderer) {
	m.rfn(r)
}

type Geometer interface {
	Geometry() geometry.Geometry
}

func (m *mesh) Geometry() geometry.Geometry {
	return m.g
}

type Materializer interface {
	Materials() []Material
	GetMaterial(int) material.Material
	AddMaterial(material.Material, int, int)
	AddGroupMaterial(material.Material, int) error
}

func (m *mesh) Materials() []Material {
	return m.materials
}

func (m *mesh) GetMaterial(idx int) material.Material {
	for _, gm := range m.materials {
		if gm.count == 0 {
			return gm.m
		}
		if gm.start >= idx && gm.start+gm.count <= idx {
			return gm.m
		}
	}
	return nil
}

func (m *mesh) AddMaterial(a material.Material, start, count int) {
	gm := Material{m, a, m.Geometry(), start, count}
	m.materials = append(m.materials, gm)
}

var InvalidGroupIdxError = xrror.Xrror("%d is an invalid group index for graphic geometry %v").Out

func (m *mesh) AddGroupMaterial(a material.Material, gidx int) error {
	geo := m.g
	if gidx < 0 || gidx >= geo.GroupCount() {
		return InvalidGroupIdxError(gidx, geo)
	}
	group := geo.GroupAt(gidx)
	m.AddMaterial(a, group.Start, group.Count)
	return nil
}

type Material struct {
	parent Mesh
	m      material.Material
	g      geometry.Geometry
	start  int
	count  int
}

func (m *Material) Shader(r Renderer) {
	pr := r.GenerateProfile(m.m)
	r.SetProgram(r, pr)
}

func (m *Material) Render(r Renderer) {
	// establish shader to use
	m.Shader(r)

	// setup underlying material
	m.m.Provide(r)

	// setup associated geometry
	gg := m.g
	gg.Provide(r)

	// setup parent mesh
	parent := m.parent
	parent.Render(r)

	//draw
	count := m.count

	indices := gg.Indices()
	mode := parent.Mode()
	if indices.Size() > 0 {
		if count == 0 {
			count = indices.Size()
		}
		val := 4 * uint32(m.start)
		r.DrawElements(mode, int32(count), graphics.UNSIGNED_INT, r.Ptr(&val))
	} else {
		if count == 0 {
			count = gg.VBOItems()
		}
		r.DrawArrays(mode, int32(m.start), int32(count))
	}
}
