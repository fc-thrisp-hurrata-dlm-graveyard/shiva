package scene

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/Laughs-In-Flowers/shiva/lib/math"
	"github.com/Laughs-In-Flowers/shiva/lib/render"

	l "github.com/yuin/gopher-lua"
)

const lTranslateNodeClass = "NTRANSLATE"

type translateNode struct {
	*node
	*graphics.Translate
}

func Translate(tag string, v math.Vector) Node {
	return &translateNode{
		newNode(tag, func(r render.Renderer, n Node) {
			//
		}, defaultRemovalFn, defaultReplaceFn, lTranslateNodeClass, lNodeClass),
		graphics.NewTranslate(v),
	}
}

func ltranslate(L *l.LState) int {
	tagFn := tagFnFor("translate", 1)
	tag := tagFn(L)
	vec := math.UnpackToVec(L, 2, math.VEC3, true)
	t := Translate(tag, vec)
	return pushNode(L, t)
}

var lTranslateNodeTable = &lua.Table{
	lTranslateNodeClass,
	[]*lua.Table{nodeTable},
	nil, nil, nil,
}

const lScaleNodeClass = "NSCALE"

type scaleNode struct {
	*node
	*graphics.Scale
}

func Scale(tag string, v math.Vector) Node {
	return &scaleNode{
		newNode(tag, func(r render.Renderer, n Node) {
			//
		}, defaultRemovalFn, defaultReplaceFn, lScaleNodeClass, lNodeClass),
		graphics.NewScale(v),
	}
}

func lscale(L *l.LState) int {
	tagFn := tagFnFor("scale", 1)
	tag := tagFn(L)
	vec := math.UnpackToVec(L, 2, math.VEC3, true)
	t := Scale(tag, vec)
	return pushNode(L, t)
}

var lScaleNodeTable = &lua.Table{
	lScaleNodeClass,
	[]*lua.Table{nodeTable},
	nil, nil, nil,
}

const lRotateNodeClass = "NROTATE"

type rotateNode struct {
	*node
	*graphics.Rotate
}

func Rotate(tag string, q math.Quaternion) Node {
	return &rotateNode{
		newNode(tag, func(r render.Renderer, n Node) {
			//
		}, defaultRemovalFn, defaultReplaceFn, lRotateNodeClass, lNodeClass),
		graphics.NewRotate(q),
	}
}

func lrotate(L *l.LState) int {
	tagFn := tagFnFor("rotate", 1)
	tag := tagFn(L)
	//quat := nil //math.quat //math.UnpackToVec(L, 2, math.VEC3, true)
	t := Rotate(tag, nil)
	return pushNode(L, t)
}

var lRotateNodeTable = &lua.Table{
	lRotateNodeClass,
	nil, nil, nil, nil,
}
