package math

import (
	"strings"

	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/davecgh/go-spew/spew"

	glm "math"

	l "github.com/yuin/gopher-lua"
)

type Vector interface {
	MathType
	Manipulator
	Getter
	Setter
	Length
	Add(Vector) Vector
	Clone() Vector
	Cross(Vector) Vector
	Dot(Vector) float32
	Mul(float32) Vector
	Normalize() Vector
	Rotate(Quaternion) Vector
	Sub(Vector) Vector
}

type VecN struct {
	tag string
	v   []float32
}

func resolveTag(tag string, n int) string {
	if tag == "" {
		switch n {
		case 2:
			return VEC2
		case 3:
			return VEC3
		case 4:
			return VEC4
		}
		return "UNKNOWN"
	}
	return tag
}

func NewVecN(tag string, n int) *VecN {
	if shouldPool {
		return &VecN{tag: tag, v: grabFromPool(n)}
	} else {
		return &VecN{tag: tag, v: make([]float32, n)}
	}
}

func NewVecNFrom(tag string, initial []float32) *VecN {
	if initial == nil {
		return &VecN{tag: tag}
	}
	var internal []float32
	if shouldPool {
		internal = grabFromPool(len(initial))
	} else {
		internal = make([]float32, len(initial))
	}
	copy(internal, initial)
	return &VecN{tag: tag, v: internal}
}

func (v *VecN) Tag() string {
	return resolveTag(v.tag, len(v.v))
}

func (v *VecN) Raw() []float32 {
	return v.v
}

func (v *VecN) Get(i int) float32 {
	return v.v[i]
}

func oneOf(k string, o ...string) bool {
	for _, v := range o {
		if k == v {
			return true
		}
	}
	return false
}

func (v *VecN) GetStr(k string) float32 {
	var ret float32
	key := strings.ToLower(k)
	ln := len(v.v)
	switch {
	case oneOf(key, "x", "r", "s") && ln >= 1:
		ret = v.v[0]
	case oneOf(key, "y", "g", "t") && ln >= 2:
		ret = v.v[1]
	case oneOf(key, "z", "b", "p") && ln >= 3:
		ret = v.v[2]
	case oneOf(key, "w", "a", "q") && ln >= 4:
		ret = v.v[3]
	}
	return ret
}

func (v *VecN) Set(i int, val float32) {
	v.v[i] = val
}

func (v *VecN) SetStr(k string, val float32) {
	key := strings.ToLower(k)
	ln := len(v.v)
	switch {
	case oneOf(key, "x", "r", "s") && ln >= 1:
		v.v[0] = val
	case oneOf(key, "y", "g", "t") && ln >= 2:
		v.v[1] = val
	case oneOf(key, "z", "b", "p") && ln >= 3:
		v.v[2] = val
	case oneOf(key, "w", "a", "q") && ln >= 4:
		v.v[3] = val
	}
}

func (v *VecN) Update(in ...float32) {
	ln := len(v.v)
	for idx, num := range in {
		if idx <= ln {
			v.v[idx] = num
		}
	}
}

func (v *VecN) RawLen() int {
	return len(v.v)
}

func (v *VecN) Len() float32 {
	if v == nil {
		return float32(glm.NaN())
	}
	if len(v.v) == 0 {
		return 0
	}

	return float32(glm.Sqrt(float64(v.Dot(v))))
}

func (v *VecN) Clone() Vector {
	cv := NewVecN(v.tag, v.RawLen())
	for idx, val := range v.v {
		cv.Set(idx, val)
	}
	return cv
}

func (v *VecN) Normalize() Vector {
	if v == nil {
		return nil
	}
	return v.Mul(1 / v.Len())
}

func (v *VecN) resize(n int) *VecN {
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
	return v
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

func (v *VecN) Add(o Vector) Vector {
	if v == nil || o == nil {
		return nil
	}
	size := intMin(len(v.v), o.RawLen())
	dst := NewVecN("", size)

	for i := 0; i < size; i++ {
		dst.v[i] = v.v[i] + o.Get(i)
	}

	return dst
}

func (v *VecN) Sub(o Vector) Vector {
	if v == nil || o == nil {
		return nil
	}
	size := intMin(len(v.v), o.RawLen())
	dst := NewVecN("", size)

	for i := 0; i < size; i++ {
		dst.v[i] = v.v[i] - o.Get(i)
	}

	return dst
}

func (v *VecN) Mul(c float32) Vector {
	if v == nil {
		return nil
	}

	dst := NewVecN("", len(v.v))

	for i, el := range v.v {
		dst.v[i] = el * c
	}

	return dst
}

func (v *VecN) Dot(o Vector) float32 {
	if v == nil || o == nil || len(v.v) != o.RawLen() {
		return float32(glm.NaN())
	}

	var result float32 = 0.0
	for i, el := range v.v {
		result += el * o.Get(i)
	}

	return result
}

func (v *VecN) Cross(o Vector) Vector {
	if v == nil || o == nil {
		return nil
	}
	if len(v.v) != 3 || o.RawLen() != 3 {
		panic("Cannot take binary cross product of non-3D elements (7D cross product not implemented)")
	}

	dst := NewVecN(VEC3, 3)
	dst.v[0], dst.v[1], dst.v[2] = v.v[1]*o.Get(2)-v.v[2]*o.Get(1),
		v.v[2]*o.Get(0)-v.v[0]*o.Get(2),
		v.v[0]*o.Get(1)-v.v[1]*o.Get(0)
	return dst
}

func CrossVectors(a, b Vector) Vector {
	return a.Cross(b)
}

func (v *VecN) Rotate(q Quaternion) Vector {
	vx := v.Get(0)
	vy := v.Get(1)
	vz := v.Get(2)

	qw := q.Get(0)
	qx := q.Get(1)
	qy := q.Get(2)
	qz := q.Get(3)

	// calculate quat * vector
	ix := qw*vx + qy*vz - qz*vy
	iy := qw*vy + qz*vx - qx*vz
	iz := qw*vz + qx*vy - qy*vx
	iw := -qx*vx - qy*vy - qz*vz
	// calculate result * inverse quat
	v.Set(0, ix*qw+iw*-qx+iy*-qz-iz*-qy)
	v.Set(1, iy*qw+iw*-qy+iz*-qx-ix*-qz)
	v.Set(2, iz*qw+iw*-qz+ix*-qy-iy*-qx)
	return v
}

func SetVectorFromMatrice(v Vector, m Matrice, ps ...MxPos) Vector {
	for _, p := range ps {
		v.Set(p.Correspondence, m.Get(p.Row, p.Column))
	}
	return v
}

func SetVectorFromRotationMatrix(v Vector, m Matrice) Vector {
	m11 := m.Get(1, 1)
	m12 := m.Get(1, 2)
	m13 := m.Get(1, 3)
	m22 := m.Get(2, 2)
	m23 := m.Get(2, 3)
	m32 := m.Get(3, 2)
	m33 := m.Get(3, 3)

	var vx, vy, vz float32
	vy = Asin(Clamp(m13, -1, 1))
	if Abs(m13) < 0.99999 {
		vx = Atan2(-m23, m33)
		vz = Atan2(-m12, m11)
	} else {
		vx = Atan2(m32, m22)
		vz = 0
	}
	v.Update(vx, vy, vz)
	return v
}

func SetVectorFromQuaternion(v Vector, q Quaternion) Vector {
	m := Mat4()
	m.Rotate(q)
	return SetVectorFromRotationMatrix(v, m)
}

func ExpectedVector(k string, v interface{}) (Vector, bool) {
	if vec, ok := v.(Vector); ok {
		switch vec.Tag() {
		case VEC2:
			if k == VEC2 {
				return vec, true
			}
		case VEC3:
			if k == VEC3 {
				return vec, true
			}
		case VEC4:
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

func Vec2(x, y float32) *VecN {
	return NewVecNFrom(VEC2, []float32{x, y})
}

func Vec3(x, y, z float32) *VecN {
	return NewVecNFrom(VEC3, []float32{x, y, z})
}

func Vec4(x, y, z, w float32) *VecN {
	return NewVecNFrom(VEC4, []float32{x, y, z, w})
}

func VecUnp(v ...float32) *VecN {
	return NewVecNFrom("", v)
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
			case Vector:
				switch dt.Tag() {
				case VEC2:
					list = unpVec(2, dt)
				case VEC3:
					list = unpVec(3, dt)
				case VEC4:
					list = unpVec(4, dt)
				}
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
	return NewVecNFrom(k, list)
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
