package scene

import (
	"strings"

	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/Laughs-In-Flowers/shiva/lib/math"
	"github.com/Laughs-In-Flowers/shiva/lib/render"

	l "github.com/yuin/gopher-lua"
)

type TransformT int

const (
	UNKNOWN_TRANSFORM TransformT = iota
	TRANSLATE
	SCALE
	ROTATE
	DIRECTION
)

func (t TransformT) String() string {
	switch t {
	case TRANSLATE:
		return "translate"
	case SCALE:
		return "scale"
	case ROTATE:
		return "rotate"
	case DIRECTION:
		return "direction"
	}
	return "unknown"
}

func stringToTransform(s string) TransformT {
	switch strings.ToLower(s) {
	case "translate":
		return TRANSLATE
	case "scale":
		return SCALE
	case "rotate":
		return ROTATE
	case "direction":
		return DIRECTION
	}
	return UNKNOWN_TRANSFORM
}

type Transform interface {
	Raw() []float32
	Set(int, float32)
	SetStr(string, float32)
	Update(...float32)
}

type translate struct {
	math.Vector
}

func newTranslate(v math.Vector) *translate {
	if v == nil {
		v = math.Vec3(0, 0, 0)
	}
	return &translate{v}
}

const lTranslateNodeClass = "NTRANSLATE"

type translateNode struct {
	*translate
	*node
}

func Translate(tag string, vec math.Vector) Node {
	return &translateNode{
		newTranslate(vec),
		newNode(tag, func(r render.Renderer, n Node) {
			//
		}, defaultRemovalFn, defaultReplaceFn, lTranslateNodeClass, lNodeClass),
	}
}

var translateTag TagFunc = tagFnFor("translate", 1)

func ltranslate(L *l.LState) int {
	tag := translateTag(L)
	vec := math.UnpackToVec(L, 2, math.VEC3, true)
	t := Translate(tag, vec)
	return pushNode(L, t)
}

var lTranslateNodeTable = &lua.Table{
	lTranslateNodeClass,
	[]*lua.Table{nodeTable},
	nil, nil, nil,
}

type scale struct {
	math.Vector
}

func newScale(v math.Vector) *scale {
	if v == nil {
		v = math.Vec3(1, 1, 1)
	}
	return &scale{v}
}

const lScaleNodeClass = "NSCALE"

type scaleNode struct {
	*scale
	*node
}

func Scale(tag string, vec math.Vector) Node {
	return &scaleNode{
		newScale(vec),
		newNode(tag, func(r render.Renderer, n Node) {
			//
		}, defaultRemovalFn, defaultReplaceFn, lScaleNodeClass, lNodeClass),
	}
}

var scaleTag TagFunc = tagFnFor("scale", 1)

func lscale(L *l.LState) int {
	tag := scaleTag(L)
	vec := math.UnpackToVec(L, 2, math.VEC3, true)
	t := Scale(tag, vec)
	return pushNode(L, t)
}

var lScaleNodeTable = &lua.Table{
	lScaleNodeClass,
	[]*lua.Table{nodeTable},
	nil, nil, nil,
}

type rotate struct {
	math.Quaternion
}

func newRotate(q math.Quaternion) *rotate {
	if q == nil {
		q = math.Quat(0, 0, 0, 1)
	}
	return &rotate{q}
}

const lRotateNodeClass = "NROTATE"

type rotateNode struct {
	*rotate
	*node
}

func Rotate(tag string, q math.Quaternion) Node {
	return &rotateNode{
		newRotate(q),
		newNode(tag, func(r render.Renderer, n Node) {
			//
		}, defaultRemovalFn, defaultReplaceFn, lRotateNodeClass, lNodeClass),
	}
}

var rotateTag TagFunc = tagFnFor("rotate", 1)

func lrotate(L *l.LState) int {
	tag := rotateTag(L)
	q := math.UnpackToQuat(L, 2, true)
	t := Rotate(tag, q)
	return pushNode(L, t)
}

var lRotateNodeTable = &lua.Table{
	lRotateNodeClass,
	[]*lua.Table{nodeTable},
	nil, nil, nil,
}

type direction struct {
	math.Vector
}

func newDirection(v math.Vector) *direction {
	if v == nil {
		v = math.Vec3(0, 0, 1)
	}
	return &direction{v}
}

const lDirectionNodeClass = "NDIRECTION"

type directionNode struct {
	*direction
	*node
}

func Direction(tag string, vec math.Vector) Node {
	return &directionNode{
		newDirection(vec),
		newNode(tag, func(r render.Renderer, n Node) {
			//
		}, defaultRemovalFn, defaultReplaceFn, lDirectionNodeClass, lNodeClass),
	}
}

var directionTag TagFunc = tagFnFor("direction", 1)

func ldirection(L *l.LState) int {
	tag := directionTag(L)
	vec := math.UnpackToVec(L, 2, math.VEC3, true)
	t := Direction(tag, vec)
	return pushNode(L, t)
}

var lDirectionNodeTable = &lua.Table{
	lDirectionNodeClass,
	[]*lua.Table{nodeTable},
	nil, nil, nil,
}

type position struct {
	t, s, r, d  Transform
	matrix      math.Matrice
	matrixWorld math.Matrice
}

func newPosition() *position {
	return &position{
		newTranslate(math.Vec3(0, 0, 0)),
		newScale(math.Vec3(0, 0, 0)),
		newRotate(math.Quat(0, 0, 0, 0)),
		newDirection(math.Vec3(0, 0, 0)),
		math.IdentityMatrix(math.MAT4),
		math.IdentityMatrix(math.MAT4),
	}
}

func (o *position) getTransform(key TransformT) Transform {
	var u Transform
	switch key {
	case TRANSLATE:
		u = o.t
	case SCALE:
		u = o.s
	case ROTATE:
		u = o.r
	case DIRECTION:
		u = o.d
	}
	return u
}

func (o *position) Set(key TransformT, sub string, val float32) {
	u := o.getTransform(key)
	u.SetStr(sub, val)
}

func (o *position) Update(key TransformT, val ...float32) {
	u := o.getTransform(key)
	u.Update(val...)
}

func (o *position) updateMatrix() {
	o.matrix.Compose(o.t.Raw(), o.r.Raw(), o.s.Raw())
}

func (o *position) updateMatrixWorld(r render.Renderer) {
	o.updateMatrix()
	last := r.Last()
	o.matrixWorld = math.MultiplyMatrices(last, o.matrix)
	r.SetLast(o.matrixWorld)
}

func (o *position) translate(world math.Vector) math.Vector {
	return math.SetVectorFromMatrice(world, o.matrixWorld, math.TranslateMxPos...)
}

func (o *position) scale(world math.Vector) math.Vector {
	o.matrixWorld.Decompose(nil, world, nil)
	return world
}

func (o *position) rotate(world math.Quaternion) math.Quaternion {
	o.matrixWorld.Decompose(nil, nil, world)
	return world
}

func (o *position) direct(world math.Vector) math.Vector {
	var q math.Quaternion
	q = math.Quat(0, 0, 0, 0)
	q = o.rotate(q)
	world.Update(o.d.Raw()...)
	world.Rotate(q)
	return world
}

type positionNode struct {
	*position
	*node
}

const lPositionNodeClass = "NPOSITION"

func Position(tag string) Node {
	pp := newPosition()
	nn := newNode(tag, func(r render.Renderer, n Node) {
		//
		pp.updateMatrixWorld(r)
	}, defaultRemovalFn, defaultReplaceFn, lPositionNodeClass, lNodeClass)

	return &positionNode{
		pp,
		nn,
	}
}

var positionTag TagFunc = tagFnFor("postion", 1)

func lposition(L *l.LState) int {
	tag := directionTag(L)
	p := Position(tag)
	return pushNode(L, p)
}

var lPositionNodeTable = &lua.Table{
	lPositionNodeClass,
	[]*lua.Table{nodeTable},
	nil, nil, nil,
}
