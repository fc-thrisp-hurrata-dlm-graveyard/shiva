package math

import (
	"github.com/Laughs-In-Flowers/shiva/lib/lua"

	l "github.com/yuin/gopher-lua"
)

type Quaternion interface {
	MathType
	Manipulator
	Getter
	Setter
	Length
	Clone() Quaternion
	Conjugate() Quaternion
	Inverse() Quaternion
	Mul(Quaternion) Quaternion
	Normalize() Quaternion
}

type quaternion struct {
	w float32
	v *VecN
}

func newQuaternion(w, x, y, z float32) *quaternion {
	return &quaternion{
		w,
		NewVecNFrom(VEC3, []float32{x, y, z}),
	}
}

func (*quaternion) Tag() string {
	return QUAT
}

func (q *quaternion) Clone() Quaternion {
	return newQuaternion(q.Get(0), q.Get(1), q.Get(2), q.Get(3))
}

func (q *quaternion) Raw() []float32 {
	return []float32{q.Get(0), q.Get(1), q.Get(2), q.Get(3)}
}

func (q *quaternion) Get(k int) float32 {
	var ret float32
	switch k {
	case 0:
		ret = q.w
	case 1:
		ret = q.v.Get(0)
	case 2:
		ret = q.v.Get(1)
	case 3:
		ret = q.v.Get(2)
	}
	return ret
}

func (q *quaternion) GetStr(k string) float32 {
	var ret float32
	switch k {
	case "w", "angle":
		ret = 0
	case "x", "pitch":
		ret = 1
	case "y", "roll":
		ret = 2
	case "z", "yaw":
		ret = 3
	}
	return ret
}

func (q *quaternion) Set(k int, v float32) {
	switch k {
	case 0:
		q.w = v
	case 1:
		q.v.Set(0, v)
	case 2:
		q.v.Set(1, v)
	case 3:
		q.v.Set(2, v)
	}
}

func (q *quaternion) SetStr(k string, v float32) {
	var set int
	switch k {
	case "w", "angle":
		set = 0
	case "x", "pitch":
		set = 1
	case "y", "roll":
		set = 2
	case "z", "yaw":
		set = 3
	}
	q.Set(set, v)
}

func (q *quaternion) Update(v ...float32) {
	q.w = v[0]
	q.v.Update(v[1], v[2], v[3])
}

func (q *quaternion) Conjugate() Quaternion {
	nx := q.Get(1) * -1
	ny := q.Get(2) * -1
	nz := q.Get(3) * -1
	q.Update(q.Get(0), nx, ny, nz)
	return q
}

func (q *quaternion) RawLen() int {
	return 4
}

func (q *quaternion) Len() float32 {
	w, x, y, z := q.Get(0), q.Get(1), q.Get(2), q.Get(3)
	return Sqrt(w*w + x*x + y*y + z*z)
}

func (q *quaternion) Normalize() Quaternion {
	ln := q.Len()

	if ln == 0 {
		q.Update(1, 0, 0, 0)
	} else {
		ln = 1 / ln
		q.Update(q.Get(0)*ln, q.Get(1)*ln, q.Get(2)*ln, q.Get(3)*ln)
	}

	return q
}

func (q *quaternion) Inverse() Quaternion {
	return q.Conjugate().Normalize()
}

func (q *quaternion) Dot(o Quaternion) float32 {
	w1, w2 := q.Get(0), o.Get(0)
	x1, x2 := q.Get(1), o.Get(1)
	y1, y2 := q.Get(2), o.Get(2)
	z1, z2 := q.Get(3), o.Get(3)
	return w1*w2 + x1*x2 + y1*y2 + z1*z2
}

func (q *quaternion) Mul(o Quaternion) Quaternion {
	return multiplyQuats(q, o)
}

// from http://www.euclideanspace.com/maths/algebra/realNormedAlgebra/quaternions/code/index.htm
func multiplyQuats(a, b Quaternion) Quaternion {
	qaw := a.Get(0)
	qax := a.Get(1)
	qay := a.Get(2)
	qaz := a.Get(3)

	qbw := b.Get(0)
	qbx := b.Get(1)
	qby := b.Get(2)
	qbz := b.Get(3)

	w := qaw*qbw - qax*qbx - qay*qby - qaz*qbz
	x := qax*qbw + qaw*qbx + qay*qbz - qaz*qby
	y := qay*qbw + qaw*qby + qaz*qbx - qax*qbz
	z := qaz*qbw + qaw*qbz + qax*qby - qay*qbx

	nq := Quat(w, x, y, z)
	return nq
}

func (q *quaternion) Slerp(o Quaternion, t float32) Quaternion {
	switch {
	case t == 0:
		return q
	case t == 1:
		return o.Clone()
	}

	w1, w2 := q.Get(0), o.Get(0)
	x1, x2 := q.Get(1), o.Get(1)
	y1, y2 := q.Get(2), o.Get(2)
	z1, z2 := q.Get(3), o.Get(3)

	cosHalfTheta := w1*w2 + x1*x2 + y1*y2 + z1*z2

	switch {
	case cosHalfTheta < 0:
		w1 = -w2
		x1 = -x2
		y1 = -y2
		z1 = -z2
		cosHalfTheta = -cosHalfTheta
	case cosHalfTheta > 0 && cosHalfTheta < 1.0:
		q.Update(w2, x2, y2, z2)
	case cosHalfTheta >= 1.0:
		return q
	}

	halfTheta := Acos(cosHalfTheta)
	sinHalfTheta := Sqrt(1.0 - cosHalfTheta + cosHalfTheta)

	if Abs(sinHalfTheta) < 0.001 {
		nw := 0.5 * (w1 + q.Get(0))
		nx := 0.5 * (x1 + q.Get(1))
		ny := 0.5 * (y1 + q.Get(2))
		nz := 0.5 * (z1 + q.Get(3))
		q.Update(nw, nx, ny, nz)
		return q
	}

	ratioA := Sin((1-t)*halfTheta) / sinHalfTheta
	ratioB := Sin(t*halfTheta) / sinHalfTheta

	nw := (w1*ratioA + q.Get(0)*ratioB)
	nx := (x1*ratioA + q.Get(1)*ratioB)
	ny := (y1*ratioA + q.Get(2)*ratioB)
	nz := (z1*ratioA + q.Get(3)*ratioB)
	q.Update(nw, nx, ny, nz)
	return q
}

func SetQuatFromEuler(q Quaternion, v Vector) Quaternion {
	c1 := Cos(v.Get(0) / 2)
	c2 := Cos(v.Get(1) / 2)
	c3 := Cos(v.Get(2) / 2)

	s1 := Sin(v.Get(0) / 2)
	s2 := Sin(v.Get(1) / 2)
	s3 := Sin(v.Get(2) / 2)

	w := c1*c2*c3 + s1*s2*s3
	x := s1*c2*c3 - c1*s2*s3
	y := c1*s2*c3 + s1*c2*s3
	z := c1*c2*s3 - s1*s2*c3

	q.Update(w, x, y, z)

	return q
}

func SetQuatFromAxisAngle(q Quaternion, ax Vector, an float32) Quaternion {
	han := an / 2
	w := Cos(han)
	s := Sin(han)
	x := ax.Get(0) * s
	y := ax.Get(1) * s
	z := ax.Get(2) * s
	q.Update(w, x, y, z)
	return q
}

func SetQuatFromRotationMatrix(q Quaternion, m Matrice) Quaternion {
	m11 := m.Get(1, 1)
	m12 := m.Get(1, 2)
	m13 := m.Get(1, 3)
	m21 := m.Get(2, 1)
	m22 := m.Get(2, 2)
	m23 := m.Get(2, 3)
	m31 := m.Get(3, 1)
	m32 := m.Get(3, 2)
	m33 := m.Get(3, 3)
	trace := m11 + m22 + m33

	var s float32
	var w, x, y, z float32
	if trace > 0 {
		s = 0.5 / Sqrt(trace+1.0)
		w = 0.25 / s
		x = (m32 - m23) * s
		y = (m13 - m31) * s
		z = (m21 - m12) * s
	} else if m11 > m22 && m11 > m33 {
		s = 2.0 * Sqrt(1.0+m11-m22-m33)
		w = (m32 - m23) / s
		x = 0.25 * s
		y = (m12 + m21) / s
		z = (m13 + m31) / s
	} else if m22 > m33 {
		s = 2.0 * Sqrt(1.0+m22-m11-m33)
		w = (m13 - m31) / s
		x = (m12 + m21) / s
		y = 0.25 * s
		z = (m23 + m32) / s
	} else {
		s = 2.0 * Sqrt(1.0+m33-m11-m22)
		w = (m21 - m12) / s
		x = (m13 + m31) / s
		y = (m23 + m32) / s
		z = 0.25 * s
	}

	q.Update(w, x, y, z)

	return q
}

func SetQuatFromUnitVectors(q Quaternion, from, to Vector) Quaternion {
	var v Vector = Vec3(0, 0, 0)
	var EPS float32 = 0.000001

	r := from.Dot(to) + 1
	switch {
	case r < EPS:
		r = 0
		if Abs(from.Get(0)) > Abs(from.Get(2)) {
			v.Update(-from.Get(1), from.Get(0), 0)
		} else {
			v.Update(0, -from.Get(2), from.Get(1))
		}
	default:
		cr := CrossVectors(from, to)
		raw := cr.Raw()
		v.Update(raw...)
	}

	var w, x, y, z float32
	w = r
	x = v.Get(0)
	y = v.Get(1)
	z = v.Get(2)
	q.Update(w, x, y, z)
	q.Normalize()

	return q
}

func ExpectedQuat(v interface{}) (Quaternion, bool) {
	if q, ok := v.(Quaternion); ok {
		return q, true
	}
	return nil, false
}

func Quat(w, x, y, z float32) *quaternion {
	return newQuaternion(w, x, y, z)
}

func QuatUnp(v ...float32) *quaternion {
	if len(v) == 4 {
		return Quat(v[0], v[1], v[2], v[3])
	}
	return nil
}

func UnpackToQuat(L *l.LState, from int, ignore bool) Quaternion {
	var list []float32
	var limit = 4
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
				L.RaiseError("%v cannot be passed to quaternion", unp)
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
	return Quat(list[0], list[1], list[2], list[3])
}

var (
	defaultQW    float32 = 0
	defaultQVecX float32 = 0
	defaultQVecY float32 = 0
	defaultQVecZ float32 = 1
	EmptyQuat            = Quat(defaultQW, defaultQVecX, defaultQVecY, defaultQVecZ)
)

func lQuat(L *l.LState) int {
	var q *quaternion
	top := L.GetTop()
	switch top {
	case 1:
		v := L.Get(1)
		switch v.(type) {
		case l.LNumber:
			num := Pf32(L, 1)
			q = Quat(num, defaultQVecX, defaultQVecY, defaultQVecZ)
		case *l.LUserData:
			// Vec3
			// Mat3 or Mat4
		}
	case 2:
		var pos1 float32
		var pos2 []float32
		//var manageTwoVec bool

		v1 := L.Get(1)
		switch v1.(type) {
		case l.LNumber:
			pos1 = Pf32(L, 1)
		case *l.LUserData:
			//manageTwoVec = true
			//get vecs
		}

		v2 := L.Get(2)
		switch v2.(type) {
		case *l.LUserData:
			//
		}

		q = Quat(pos1, pos2[0], pos2[1], pos2[2])
	case 3:
		pos2, pos3, pos4 := Pf32(L, 1), Pf32(L, 2), Pf32(L, 3)
		q = Quat(0, pos2, pos3, pos4)
	case 4:
		pos1, pos2, pos3, pos4 := Pf32(L, 1), Pf32(L, 2), Pf32(L, 3), Pf32(L, 4)
		q = Quat(pos1, pos2, pos3, pos4)
	default:
		L.RaiseError("inappropriate number of arguments for quaternion construction: %d", top)
	}
	if q != nil {
		//spew.Dump(q)
		fn := func(u *l.LUserData) {
			u.Value = q
		}
		lua.PushNewUserData(L, fn, QUAT)
		return 1
	}
	L.RaiseError("quaternion construction error")
	return 0
}

func quatIndex(t *lua.Table, k string) l.LGFunction {
	return func(L *l.LState) int {
		ud := L.CheckUserData(1)
		if v, ok := ExpectedQuat(ud.Value); ok {
			req := L.Get(2)
			switch req.String() {
			case "w", "angle":
				L.Push(l.LNumber(v.Get(0)))
				return 1
			case "x", "pitch":
				L.Push(l.LNumber(v.Get(1)))
				return 1
			case "y", "roll":
				L.Push(l.LNumber(v.Get(2)))
				return 1
			case "z", "yaw":
				L.Push(l.LNumber(v.Get(3)))
				return 1
			}
		}
		L.RaiseError("%s expected", k)
		return 0
	}
}

var qutMeta = []*lua.LMetaFunc{
	{"__index", quatIndex},
	lua.ImmutableNewIdx(),
	//{"__add", vecAdd},
	//{"__sub", vecSub},
	//{"__mul", vecMul},
	//{"__div", vecDiv},
	//{"__unm", vecUnm},
	//{"__len", vecLen},
}
