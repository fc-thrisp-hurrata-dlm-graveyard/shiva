package scene

import (
	"math"

	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/geometry"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/material"
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/Laughs-In-Flowers/shiva/lib/render"

	l "github.com/yuin/gopher-lua"
)

type sphere struct {
	*node
	m render.Mesh
}

const lSphereNodeClass = "NSPHERE"

func Sphere(tag string) *sphere {
	s := geometry.Sphere(1, 16, 16, 0, math.Pi*2, 0, math.Pi)

	sm := render.NewMesh(
		"SPHERE",
		s,
		func(r render.Renderer) {
			//log.Println("sphere render func")
			//spew.Dump("sphere mesh inner")
		},
		graphics.TRIANGLES,
	)

	mat := material.Basic()
	mat.SetWireframe(true)
	sm.AddMaterial(mat, 0, 0)

	return &sphere{
		newNode(tag, func(r render.Renderer, n Node) {
			for _, m := range sm.Materials() {
				m.Render(r)
			}
			//spew.Dump("sphere node render")
		}, defaultRemovalFn, defaultReplaceFn, lSphereNodeClass, lNodeClass),
		sm,
	}
}

func lsphere(L *l.LState) int {
	s := Sphere("test-sphere")
	return pushNode(L, s)
}

var lSphereNodeTable = &lua.Table{
	lSphereNodeClass,
	[]*lua.Table{nodeTable},
	defaultIdxMetaFuncs(),
	nil, nil,
}
