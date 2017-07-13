package scene

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/geometry"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/material"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/mesh"
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/Laughs-In-Flowers/shiva/lib/math"
	"github.com/Laughs-In-Flowers/shiva/lib/render"

	l "github.com/yuin/gopher-lua"
)

type lines struct {
	*node
	mesh.Mesh
	mvpm graphics.Uniform // UniformMatrix4fv // Model view projection matrix uniform
}

func NewLines(tag string, g geometry.Geometry, m material.Material) *lines {
	mvpm := graphics.UniformMatrix4fv("MVP")

	li := mesh.New(
		"LINES",
		g,
		func(p graphics.Provider) {
			// Calculates model view projection matrix and updates uniform
			//mw := l.MatrixWorld()
			//var mvpm math32.Matrix4
			//mvpm.MultiplyMatrices(&rinfo.ViewMatrix, &mw)
			//mvpm.MultiplyMatrices(&rinfo.ProjMatrix, &mvpm)
			//l.mvpm.SetMatrix4(&mvpm)
			//l.mvpm.Transfer(gs)
			mvpm.Transfer(p)
		},
		graphics.LINES,
	)
	li.AddMaterial(m, 0, 0)

	n := newNode(tag, func(r render.Renderer, n Node) {
		materials := li.Materials()
		for _, m := range materials {
			m.Shader(r)
			m.Provide(r)
		}
	}, defaultRemovalFn, defaultReplaceFn, lAxisNodeClass, lNodeClass)

	return &lines{n, li, mvpm}
}

const lAxisNodeClass = "NAXIS"

func Axis(tag string, size float32) *lines {
	geo := geometry.New()
	positions := math.NewAF32(0, 18)
	positions.Append(
		0, 0, 0, size, 0, 0,
		0, 0, 0, 0, size, 0,
		0, 0, 0, 0, 0, size,
	)
	colors := math.NewAF32(0, 18)
	colors.Append(
		1, 0, 0, 1, 0.6, 0,
		0, 1, 0, 0.6, 1, 0,
		0, 0, 1, 0, 0.6, 1,
	)
	geo.AddVBO(graphics.NewBuff().AddAttrib("VertexPosition", 3).SetBuffer(positions))
	geo.AddVBO(graphics.NewBuff().AddAttrib("VertexColor", 3).SetBuffer(colors))

	mat := material.Basic()
	mat.SetLineWidth(2.0)

	axis := NewLines(tag, geo, mat)

	return axis
}

func laxis(L *l.LState) int {
	size := math.Pf32(L, 1)
	a := Axis("axis", size)
	return pushNode(L, a)
}

var lAxisNodeTable = &lua.Table{
	lAxisNodeClass,
	[]*lua.Table{nodeTable},
	nil, nil, nil,
}
