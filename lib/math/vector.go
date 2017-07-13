package math

import (
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/davecgh/go-spew/spew"

	glm "math"

	l "github.com/yuin/gopher-lua"
)

type Vector interface {
	MathType
	Raw() []float32
	Get(int) float32
	Set(int, float32)
	Len() float32
	Normalize() Vector
}

func ExpectedVector(k string, v interface{}) (Vector, bool) {
	if vec, ok := v.(Vector); ok {
		switch vec.(type) {
		case *vec2:
			if k == VEC2 {
				return vec, true
			}
		case *vec3:
			if k == VEC3 {
				return vec, true
			}
		case *vec4:
			if k == VEC4 {
				return vec, true
			}
		}
	}
	return nil, false

}

func ToVector(L *l.LState, v l.LValue) Vector {
	return toVector(L, v, defaultToErrMsgs("Vector"))
}

func toVector(L *l.LState, v l.LValue, msg toErrMsgs) Vector {
	if vv, ok := v.(*l.LUserData); ok {
		val := vv.Value
		if vec, ok := val.(Vector); ok {
			return vec
		}
		L.RaiseError(msg[0], val)
	}
	L.RaiseError(msg[1], v)
	return nil
}

type VecN struct {
	v []float32
}

func NewVecN(n int) *VecN {
	if shouldPool {
		return &VecN{v: grabFromPool(n)}
	} else {
		return &VecN{v: make([]float32, n)}
	}
}

func NewVecNFrom(initial []float32) *VecN {
	if initial == nil {
		return &VecN{}
	}
	var internal []float32
	if shouldPool {
		internal = grabFromPool(len(initial))
	} else {
		internal = make([]float32, len(initial))
	}
	copy(internal, initial)
	return &VecN{v: internal}
}

func (v VecN) Raw() []float32 {
	return v.v
}

func (v VecN) Get(i int) float32 {
	return v.v[i]
}

func (v *VecN) Set(i int, val float32) {
	v.v[i] = val
}

func (v *VecN) destroy() {
	if v == nil || v.v == nil {
		return
	}

	if shouldPool {
		returnToPool(v.v)
	}
	v.v = nil
}

func (v *VecN) Len() float32 {
	if v == nil {
		return float32(glm.NaN())
	}
	if len(v.v) == 0 {
		return 0
	}

	return float32(glm.Sqrt(float64(v.dot(v))))
}

//func (v *VecN) Resize(n int) *VecN {
//	if v == nil {
//		return NewVecN(n)
//	}
//
//	if n <= cap(v.v) {
//		if v.v != nil {
//			v.v = v.v[:n]
//		} else {
//			v.v = []float32{}
//		}
//		return v
//	}
//
//	if shouldPool && v.v != nil {
//		returnToPool(v.v)
//	}
//	*v = (*NewVecN(n))
//
//	return v
//}

//func (v *VecN) SetBackingSlice(s []float32) {
//	v.v = s
//}

//func (v *VecN) Size() int {
//	return len(v.v)
//}

//func (v *VecN) Cap() int {
//	return cap(v.v)
//}

//func (v *VecN) zero(n int) {
//	v.Resize(n)
//	for i := range v.v {
//		v.v[i] = 0
//	}
//}

func (v *VecN) add(a *VecN) *VecN {
	if v == nil || a == nil {
		return nil
	}
	size := intMin(len(v.v), len(a.v))
	dst := NewVecN(size)

	for i := 0; i < size; i++ {
		dst.v[i] = v.v[i] + a.v[i]
	}

	return dst
}

func (v *VecN) sub(s *VecN) *VecN {
	if v == nil || s == nil {
		return nil
	}
	size := intMin(len(v.v), len(s.v))
	dst := NewVecN(size)

	for i := 0; i < size; i++ {
		dst.v[i] = v.v[i] - s.v[i]
	}

	return dst
}

func (v *VecN) mul(c float32) *VecN {
	if v == nil {
		return nil
	}

	dst := NewVecN(len(v.v))

	for i, el := range v.v {
		dst.v[i] = el * c
	}

	return dst
}

func (v *VecN) dot(o *VecN) float32 {
	if v == nil || o == nil || len(v.v) != len(o.v) {
		return float32(glm.NaN())
	}

	var result float32 = 0.0
	for i, el := range v.v {
		result += el * o.v[i]
	}

	return result
}

func (v *VecN) cross(o *VecN) *VecN {
	if v == nil || o == nil {
		return nil
	}
	if len(v.v) != 3 || len(o.v) != 3 {
		panic("Cannot take binary cross product of non-3D elements (7D cross product not implemented)")
	}

	dst := NewVecN(3)
	dst.v[0], dst.v[1], dst.v[2] = v.v[1]*o.v[2]-v.v[2]*o.v[1], v.v[2]*o.v[0]-v.v[0]*o.v[2], v.v[0]*o.v[1]-v.v[1]*o.v[0]
	return dst
}

func (v *VecN) Normalize() *VecN {
	if v == nil {
		return nil
	}
	return v.mul(1 / v.Len())
}

type vec2 struct {
	*VecN
}

func (*vec2) Tag() string {
	return VEC2
}

func (v *vec2) Normalize() Vector {
	return v.Normalize()
}

func Vec2(x, y float32) *vec2 {
	return &vec2{
		NewVecNFrom([]float32{x, y}),
	}
}

type vec3 struct {
	*VecN
}

func (*vec3) Tag() string {
	return VEC3
}

func (v *vec3) Normalize() Vector {
	return v.Normalize()
}

func Vec3(x, y, z float32) *vec3 {
	return &vec3{
		NewVecNFrom([]float32{x, y, z}),
	}
}

type vec4 struct {
	*VecN
}

func (*vec4) Tag() string {
	return VEC4
}

func (v *vec4) Normalize() Vector {
	return v.Normalize()
}

func Vec4(x, y, z, w float32) *vec4 {
	return &vec4{
		NewVecNFrom([]float32{x, y, z, w}),
	}
}

func lVec(k string) l.LGFunction {
	return func(L *l.LState) int {
		v := UnpackToVec(L, 1, k, false)
		fn := func(u *l.LUserData) {
			u.Value = v
		}
		lua.PushNewUserData(L, fn, k)
		return 1
	}
}

func UnpackToVec(L *l.LState, from int, k string, ignore bool) Vector {
	var list []float32
	limit := establishLimit(k)
	top := L.GetTop()
	for i := from; i <= top; i++ {
		unp := L.Get(i)
		switch unp.(type) {
		case l.LNumber:
			list = append(list, float32(L.CheckNumber(i)))
		case *l.LUserData:
			unpVec := func(l int, v Vector) []float32 {
				var ret []float32
				for ii := 1; ii <= l; ii++ {
					ret = append(ret, v.Get(ii-1))
				}
				return ret
			}
			ud := unp.(*l.LUserData)
			d := ud.Value
			switch dt := d.(type) {
			case *vec2:
				list = unpVec(2, dt)
			case *vec3:
				list = unpVec(3, dt)
			case *vec4:
				list = unpVec(4, dt)
			}
		default:
			if !ignore {
				L.RaiseError("%v cannot be passed to vector", unp)
			}
		}
	}
	ll := len(list)
	switch {
	case ll < limit:
		for i := ll; i <= limit; i++ {
			list = append(list, 0)
		}
	case ll > limit:
		for i := ll; i > limit; i-- {
			idx := i - 1
			list = append(list[:idx], list[idx+1:]...)
		}
	}
	rv := NewVecNFrom(list)
	var ret Vector
	switch k {
	case VEC2:
		ret = &vec2{rv}
	case VEC3:
		ret = &vec3{rv}
	case VEC4:
		ret = &vec4{rv}
	}
	return ret
}

func vecIndex(t *lua.Table, _ string) l.LGFunction {
	k := t.Name
	limit := establishLimit(k)
	vecLimit := func(idx, lim int) bool {
		if lim < idx {
			return true
		}
		return false
	}
	return func(L *l.LState) int {
		ud := L.CheckUserData(1)
		if v, ok := ExpectedVector(k, ud.Value); ok {
			req := L.Get(2)
			switch req.String() {
			case "x", "r", "s":
				L.Push(l.LNumber(v.Get(0)))
				return 1
			case "y", "g", "t":
				L.Push(l.LNumber(v.Get(1)))
				return 1
			case "z", "b", "p":
				if vecLimit(3, limit) {
					L.RaiseError("'%s' not accessible to %s", req.String(), k)
					return 0
				}
				L.Push(l.LNumber(v.Get(2)))
				return 1
			case "w", "a", "q":
				if vecLimit(4, limit) {
					L.RaiseError("'%s' not accessible to %s", req.String(), k)
					return 0
				}
				L.Push(l.LNumber(v.Get(3)))
				return 1
			}
		}
		L.RaiseError("%s expected", k)
		return 0
	}
}

func VecOfSize(sz int, f []float32) (Vector, error) {
	switch sz {
	case 2:
		return Vec2(f[0], f[1]), nil
	case 3:
		return Vec3(f[0], f[1], f[2]), nil
	case 4:
		return Vec4(f[0], f[1], f[2], f[3]), nil
	}
	return nil, NoXOfSizeError("vector", sz)
}

func vecIMF(L *l.LState, k string) innerMathFunc {
	return func(v1, v2 []float32, imr innerMathRun) int {
		size, dst := imr(v1, v2)
		nv, err := VecOfSize(size, dst)
		if err != nil {
			L.RaiseError("unable to perform vector %s: %s", k, err.Error())
		}
		fn := func(u *l.LUserData) {
			u.Value = nv
		}
		mt := establishKindOfLimit("vector", size)
		lua.PushNewUserData(L, fn, mt)
		return 1
	}
}

func vecVaryParams(L *l.LState) ([]float32, []float32) {
	var rv1, rv2 []float32

	var vpos1ErrMsgs = toErrMsgs{
		"first position param %s is not UserData",
		"first position param %s is not Vector",
	}

	v1 := toVector(L, L.Get(1), vpos1ErrMsgs)
	rv1 = v1.Raw()
	sz := len(rv1)

	v2 := L.Get(2)
	switch v2.(type) {
	case l.LNumber:
		c := Pf32(L, 2)
		var cv []float32
		for i := 0; i <= sz; i++ {
			cv = append(cv, c)
		}
		rv2 = cv
	case *l.LUserData:
		v2v := ToVector(L, v2)
		rv2 = v2v.Raw()
	default:
		L.RaiseError("%s not a vector applicable math type", v2)
	}

	return rv1, rv2
}

func vecInner(imr innerMathRun) inner {
	return func(L *l.LState, vm innerMathFunc) int {
		v1, v2 := vecVaryParams(L)
		return vm(v1, v2, imr)
	}
}

func vecAdd(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		imr := func(v1, v2 []float32) (int, []float32) {
			size := intMin(len(v1), len(v2))
			dst := make([]float32, size)
			for i := 1; i <= size; i++ {
				idx := i - 1
				dst[idx] = v1[idx] + v2[idx]
			}
			return size, dst
		}
		return runInner(
			L,
			vecInner(imr),
			vecIMF(L, "addition"),
		)
	}
}

func vecSub(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		imr := func(v1, v2 []float32) (int, []float32) {
			size := intMin(len(v1), len(v2))
			dst := make([]float32, size)
			for i := 1; i <= size; i++ {
				idx := i - 1
				dst[idx] = v1[idx] - v2[idx]
			}
			return size, dst
		}
		return runInner(
			L,
			vecInner(imr),
			vecIMF(L, "subtraction"),
		)
	}
}

func vecMul(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		imr := func(v1, v2 []float32) (int, []float32) {
			size := intMin(len(v1), len(v2))
			dst := make([]float32, size)
			for i := 1; i <= size; i++ {
				idx := i - 1
				dst[idx] = v1[idx] * v2[idx]
			}
			return size, dst
		}
		return runInner(
			L,
			vecInner(imr),
			vecIMF(L, "multiplication"),
		)
		return 0
	}
}

func vecDiv(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		imr := func(v1, v2 []float32) (int, []float32) {
			size := intMin(len(v1), len(v2))
			dst := make([]float32, size)
			for i := 1; i <= size; i++ {
				idx := i - 1
				dst[idx] = v1[idx] / v2[idx]
			}
			return size, dst
		}
		return runInner(
			L,
			vecInner(imr),
			vecIMF(L, "division"),
		)
		return 0
	}
}

func vecUnm(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		k := t.Name
		ud := L.CheckUserData(1)
		if v, ok := ExpectedVector(k, ud.Value); ok {
			lt := establishLimit(k)
			for i := 0; i <= lt-1; i++ {
				nv := v.Get(i)
				nv = nv - 1
				v.Set(i, nv)
			}
			L.Push(ud)
			return 1
		}
		L.RaiseError("__unm error, %s expected", k)
		return 0
	}
}

func vecLen(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		k := t.Name
		ud := L.CheckUserData(1)
		if v, ok := ExpectedVector(k, ud.Value); ok {
			vl := v.Len()
			L.Push(l.LNumber(vl))
			return 1
		}
		L.RaiseError("__len error, %s expected", k)
		return 0
	}
}

//return raw list, take index return item
func vecCall(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		spew.Dump(L.GetTop())
		return 0
	}
}

var vecMeta = []*lua.LMetaFunc{
	{"__index", vecIndex},
	lua.ImmutableNewIdx(),
	{"__add", vecAdd},
	{"__sub", vecSub},
	{"__mul", vecMul},
	{"__div", vecDiv},
	{"__unm", vecUnm},
	{"__len", vecLen},
	{"__call", vecCall},
}
