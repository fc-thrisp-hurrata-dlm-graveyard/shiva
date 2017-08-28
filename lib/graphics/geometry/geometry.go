package geometry

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/math"
)

type Geometer interface {
	graphics.Initializer
	graphics.Closer
	graphics.Providable
	graphics.RefCounter
}

type Geometry interface {
	Geometer
	Grouper
	Indicer
	VBOer
}

type geometry struct {
	p             graphics.Provider
	refcount      int
	vbos          []*graphics.Buff
	groups        []Group
	indices       math.AU32
	handleVAO     uint32
	handleIndices graphics.Buffer
	updateIndices bool
}

func New() *geometry {
	g := &geometry{}
	g.Initialize()
	return g
}

func (g *geometry) Initialize() {
	g.p = nil
	g.refcount = 1
	g.vbos = make([]*graphics.Buff, 0)
	g.groups = make([]Group, 0)
	g.handleVAO = 0
	g.handleIndices = 0
	g.updateIndices = true
}

func (g *geometry) Close() {
	if g.refcount > 1 {
		g.refcount--
		return
	}
	if g.p != nil {
		g.p.DeleteVertexArray(g.handleVAO)
		g.p.DeleteBuffer(g.handleIndices)
	}
	for i := 0; i < len(g.vbos); i++ {
		g.vbos[i].Close()
	}
	g.Initialize()
}

func (g *geometry) Renderable() bool {
	return true
}

func (g *geometry) SetRenderable(b bool) {}

func (g *geometry) Provide(p graphics.Provider) {
	if g.p == nil {
		g.handleVAO = p.GenVertexArray()
		p.BindVertexArray(g.handleVAO) // is this necessary if being done below
		g.handleIndices = p.GenBuffer()
		g.p = p
	}

	p.BindVertexArray(g.handleVAO)
	for _, vbo := range g.vbos {
		vbo.Provide(p)
	}

	if g.indices.Size() > 0 && g.updateIndices {
		p.BindBuffer(graphics.ELEMENT_ARRAY_BUFFER, g.handleIndices)
		p.BufferData(graphics.ELEMENT_ARRAY_BUFFER, g.indices.Bytes(), p.Ptr(g.indices), graphics.STATIC_DRAW)
		g.updateIndices = false
	}
}

func (g *geometry) Increment() {
	g.refcount++
}

func (g *geometry) Decrement() {
	g.refcount--
}

type Group struct {
	Start  int
	Count  int
	MatIdx int
	MatId  string
}

type Grouper interface {
	AddGroup(int, int, int) *Group
	AddGroups(...Group)
	GroupCount() int
	GroupAt(int) *Group
}

func (g *geometry) AddGroup(start, count, matIdx int) *Group {
	g.groups = append(g.groups, Group{start, count, matIdx, ""})
	return &g.groups[len(g.groups)-1]
}

func (g *geometry) AddGroups(gs ...Group) {
	for _, ng := range gs {
		g.groups = append(g.groups, ng)
	}
}

func (g *geometry) GroupCount() int {
	return len(g.groups)
}

func (g *geometry) GroupAt(idx int) *Group {
	if idx >= 0 && idx <= len(g.groups)-1 {
		return &g.groups[idx-1]
	}
	return nil
}

type Indicer interface {
	Indices() math.AU32
	SetIndices(math.AU32)
}

func (g *geometry) Indices() math.AU32 {
	return g.indices
}

func (g *geometry) SetIndices(i math.AU32) {
	g.indices = i
	//g.boundingBoxValid = false
	//g.boundingSphereValid = false
}

type VBOer interface {
	VBO(string) *graphics.Buff
	AddVBO(*graphics.Buff)
	VBOItems() int
}

func (g *geometry) VBO(attrib string) *graphics.Buff {
	for _, vbo := range g.vbos {
		if vbo.Attrib(attrib) != nil {
			return vbo
		}
	}
	return nil
}

func (g *geometry) AddVBO(vbo *graphics.Buff) {
	g.vbos = append(g.vbos, vbo)
}

func (g *geometry) VBOItems() int {
	if len(g.vbos) == 0 {
		return 0
	}
	vbo := g.vbos[0]
	if vbo.AttribCount() == 0 {
		return 0
	}
	return vbo.Buffer().Bytes() / vbo.Stride()
}

//func (g *Geometry) BoundingBox() math32.Box3 {
//	return nil
//}

//func (g *Geometry) BoundingSphere() math32.Sphere {
//	return nil
//}
