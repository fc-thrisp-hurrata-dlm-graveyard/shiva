package mesh

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/geometry"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/material"
	"github.com/Laughs-In-Flowers/shiva/lib/render"
	"github.com/Laughs-In-Flowers/shiva/lib/xrror"
)

type Mesh interface {
	graphics.Initializer
	graphics.Closer
	graphics.Moder
	graphics.Renderable
	Geometer
	Materializer
}

type ProviderFunc func(graphics.Provider)

type mesh struct {
	tag        string
	g          geometry.Geometry
	materials  []Material
	mode       graphics.Enum
	renderable bool
	pfn        ProviderFunc
}

func New(tag string, e geometry.Geometry, pfn ProviderFunc, mode graphics.Enum) *mesh {
	m := &mesh{
		tag: tag,
		pfn: pfn,
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

func (m *mesh) Provide(p graphics.Provider) {
	m.pfn(p)
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
	m      material.Material // Associated material
	g      geometry.Geometry // Geometry from parent mesh
	start  int               // Index of first element in the geometry
	count  int               // Number of elements
}

func (m *Material) Shader(r render.Renderer) {
	pr := r.GenerateProfile(m.m)
	r.SetProgram(r, pr)
}

// render
// 1. underlying material.Material
// 2. Geometry from parent Mesh
// 3. parent Mesh
// 4. elements/arrays
func (m *Material) Provide(p graphics.Provider) {
	m.m.Provide(p)

	gg := m.g
	gg.Provide(p)

	parent := m.parent
	parent.Provide(p)

	count := m.count

	indices := gg.Indices()
	mode := parent.Mode()
	if indices.Size() > 0 {
		if count == 0 {
			count = indices.Size()
		}
		p.DrawElements(mode, int32(count), graphics.UNSIGNED_INT, p.Ptr(4*uint32(m.start)))
	} else {
		if count == 0 {
			count = gg.VBOItems()
		}
		p.DrawArrays(mode, int32(m.start), int32(count))
	}
}
