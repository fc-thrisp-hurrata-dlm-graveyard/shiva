package math

import (
	"strconv"

	"github.com/Laughs-In-Flowers/shiva/lib/lua"

	l "github.com/yuin/gopher-lua"
)

type Matrice interface {
	MathType
	Manipulator
	MultiGetter
	MultiSetter
	Length
	Add(Matrice) Matrice
	Cols() int
	Compose([]float32, []float32, []float32) Matrice
	Decompose(Vector, Vector, Quaternion) Matrice
	Determinant() float32
	LookAt(Vector, Vector, Vector) Matrice
	MulMatrice(Matrice) Matrice
	MulScalar(float32) Matrice
	MulVec(Vector) Vector
	Orthographic(float32, float32, float32, float32, float32, float32) Matrice
	Perspective(float32, float32, float32, float32) Matrice
	Rotate(Quaternion) Matrice
	Rows() int
	Scale(Vector) Matrice
	Sub(Matrice) Matrice
	Transpose() Matrice
	Trace() float32
}

type MatRxC struct {
	tag  string
	r, c int
	v    []float32
}

func newMatrix(tag string, r, c int) *MatRxC {
	if shouldPool {
		return &MatRxC{tag: tag, r: r, c: c, v: grabFromPool(r * c)}
	} else {
		return &MatRxC{tag: tag, r: r, c: c, v: make([]float32, r*c)}
	}
}

func newMatrixFromData(tag string, src []float32, r, c int) *MatRxC {
	if src == nil {
		src = make([]float32, r*c)
	}
	var internal []float32
	if shouldPool {
		internal = grabFromPool(r * c)
	} else {
		internal = make([]float32, r*c)
	}
	copy(internal, src[:r*c])

	return &MatRxC{tag: tag, r: r, c: c, v: internal}
}

func (m *MatRxC) Tag() string {
	return m.tag
}

func (m *MatRxC) Raw() []float32 {
	if m == nil {
		return nil
	}

	return m.v
}

// caution: will update the underlying array of values with whatever you provide,
// contingent on index (you cannot bring in more items than already exists in array)
func (m *MatRxC) Update(v ...float32) {
	ln := len(m.v)
	mln := ln - 1
	for idx, val := range v {
		if idx <= mln {
			m.v[idx] = val
		}
	}
}

func (m *MatRxC) Get(row, col int) float32 {
	idx := col*m.r + row
	if idx <= len(m.v)-1 {
		return m.v[idx]
	}
	return 0
}

// string row and column starting from 1, e.g. string 1,1 will pull 0,0 from the
// underlying array
func (m *MatRxC) GetStr(k1, k2 string) float32 {
	a, b := intKeys(k1, k2)
	return m.Get(a, b)
}

func (m *MatRxC) Set(row, col int, val float32) {
	m.v[col*m.r+row] = val
}

// string row and column starting from 1, e.g. string 1,1,val will push 0,0,val
// to the underlying array
func (m *MatRxC) SetStr(k1, k2 string, v float32) {
	a, b := intKeys(k1, k2)
	m.Set(a, b, v)
}

func intKeys(k1, k2 string) (int, int) {
	a, _ := strconv.ParseInt(k1, 10, 0)
	b, _ := strconv.ParseInt(k2, 10, 0)
	return int(a) - 1, int(b) - 1
}

func (m *MatRxC) Rows() int {
	return m.r
}

func (m *MatRxC) Cols() int {
	return m.c
}

// length of the underlying raw array
func (m *MatRxC) RawLen() int {
	return len(m.v)
}

// same as RawLen, until determined if or how to do otherwise
func (m *MatRxC) Len() float32 {
	return float32(len(m.v))
}

func (m *MatRxC) reshape(tag string, r, c int) *MatRxC {
	if m == nil {
		return newMatrix(tag, r, c)
	}

	if r*c <= cap(m.v) {
		if m.v != nil {
			m.v = m.v[:r*c]
		} else {
			m.v = []float32{}
		}
		m.r, m.c = r, c
		return m
	}

	if shouldPool && m.v != nil {
		returnToPool(m.v)
	}
	(*m) = (*newMatrix(tag, r, c))

	return m
}

func (m *MatRxC) Add(o Matrice) Matrice {
	if m == nil || o == nil || m.r != o.Rows() || m.c != o.Cols() {
		return nil
	}

	dst := newMatrix(m.tag, m.r, m.c)

	ov := o.Raw()
	for i, el := range m.v {
		dst.v[i] = el + ov[i]
	}

	return dst
}

func (m *MatRxC) Sub(o Matrice) Matrice {
	if m == nil || o == nil || m.r != o.Rows() || m.c != o.Cols() {
		return nil
	}

	dst := newMatrix(m.tag, m.r, m.c)

	ov := o.Raw()
	for i, el := range m.v {
		dst.v[i] = el - ov[i]
	}

	return dst
}

func (m *MatRxC) MulScalar(c float32) Matrice {
	if m == nil {
		return nil
	}

	dst := newMatrix(m.tag, m.r, m.c)

	for i, el := range m.v {
		dst.v[i] = el * c
	}

	return dst
}

func (m *MatRxC) MulVec(o Vector) Vector {
	if m == nil || o == nil || m.c != o.RawLen() {
		return nil
	}

	dst := NewVecN("", m.r)

	ov := o.Raw()
	for r := 0; r < m.r; r++ {
		dst.v[r] = 0
		for c := 0; c < m.c; c++ {
			dst.v[r] += m.Get(r, c) * ov[c]
		}
	}

	return dst
}

func MultiplyMatriceVector(m Matrice, v Vector) Vector {
	return m.MulVec(v)
}

func (m *MatRxC) MulMatrice(o Matrice) Matrice {
	or := o.Rows()
	if m == nil || o == nil || m.c != or {
		return nil
	}

	oc := o.Cols()
	dst := newMatrix(m.tag, m.r, oc)

	ov := o.Raw()
	for r1 := 0; r1 < m.r; r1++ {
		for c2 := 0; c2 < oc; c2++ {
			dst.v[c2*m.r+r1] = 0
			for i := 0; i < m.c; i++ {
				dst.v[c2*m.r+r1] += m.v[i*m.r+r1] * ov[c2*or+i]
			}

		}
	}

	return dst
}

func MultiplyMatrices(a, b Matrice) Matrice {
	return a.MulMatrice(b)
}

// TODO: rewrite to perform algorithmically based on len of underlying array values
// or other available details
func (m *MatRxC) Determinant() float32 {
	r := m.Raw()
	var res float32
	switch m.tag {
	case MAT2:
		res = r[0]*r[3] - r[1]*r[2]
	case MAT3:
		res = r[0]*r[4]*r[8] + r[3]*r[7]*r[2] + r[6]*r[1]*r[5] - r[6]*r[4]*r[2] - r[3]*r[1]*r[8] - r[0]*r[7]*r[5]
	case MAT4:
		res = r[0]*r[5]*r[10]*r[15] - r[0]*r[5]*r[11]*r[14] - r[0]*r[6]*r[9]*r[15] + r[0]*r[6]*r[11]*r[13] + r[0]*r[7]*r[9]*r[14] - r[0]*r[7]*r[10]*r[13] - r[1]*r[4]*r[10]*r[15] + r[1]*r[4]*r[11]*r[14] + r[1]*r[6]*r[8]*r[15] - r[1]*r[6]*r[11]*r[12] - r[1]*r[7]*r[8]*r[14] + r[1]*r[7]*r[10]*r[12] + r[2]*r[4]*r[9]*r[15] - r[2]*r[4]*r[11]*r[13] - r[2]*r[5]*r[8]*r[15] + r[2]*r[5]*r[11]*r[12] + r[2]*r[7]*r[8]*r[13] - r[2]*r[7]*r[9]*r[12] - r[3]*r[4]*r[9]*r[14] + r[3]*r[4]*r[10]*r[13] + r[3]*r[5]*r[8]*r[14] - r[3]*r[5]*r[10]*r[12] - r[3]*r[6]*r[8]*r[13] + r[3]*r[6]*r[9]*r[12]
	}
	return res
}

func (m *MatRxC) Transpose() Matrice {
	if m == nil {
		return nil
	}

	dst := newMatrix("", m.c, m.r)

	for r := 0; r < m.r; r++ {
		for c := 0; c < m.c; c++ {
			dst.v[r*dst.r+c] = m.v[c*m.r+r]
		}
	}

	return dst
}

func (m *MatRxC) Trace() float32 {
	if m == nil || m.r != m.c {
		return NaN()
	}

	var out float32
	for i := 0; i < m.r; i++ {
		out += m.Get(i, i)
	}

	return out
}

func SetMatriceFromVector(m Matrice, v Vector, ps ...MxPos) Vector {
	for _, p := range ps {
		m.Set(p.Row, p.Column, v.Get(p.Correspondence))
	}
	return v
}

func (m *MatRxC) Translate(v Vector) Matrice {
	SetMatriceFromVector(m, v, TranslateMxPos...)
	return m
}

func (m *MatRxC) Scale(v Vector) Matrice {
	sc := v.RawLen()
	if sc <= m.c {
		for c := 0; c <= sc-1; c++ {
			cv := v.Get(c)
			for r := 0; r <= m.r-1; r++ {
				m.Set(r, c, cv*m.Get(r, c))
			}
		}
		return m
	}
	return nil
}

func (m *MatRxC) Rotate(q Quaternion) Matrice {
	if q != nil {
		w := q.Get(0)
		x := q.Get(1)
		y := q.Get(2)
		z := q.Get(3)

		x2 := x + x
		y2 := y + y
		z2 := z + z
		xx := x * x2
		xy := x * y2
		xz := x * z2
		yy := y * y2
		yz := y * z2
		zz := z * z2
		wx := w * x2
		wy := w * y2
		wz := w * z2

		m.Set(0, 0, 1-(yy+zz))
		m.Set(0, 1, xy-wz)
		m.Set(0, 2, xz+wy)

		m.Set(1, 0, xy+wz)
		m.Set(1, 1, 1-(xx+zz))
		m.Set(1, 2, yz-wx)

		m.Set(2, 0, xz-wy)
		m.Set(2, 1, yz+wx)
		m.Set(2, 2, 1-(xx+yy))

		m.Set(3, 0, 0)
		m.Set(3, 1, 0)
		m.Set(3, 2, 0)

		m.Set(0, 3, 0)
		m.Set(1, 3, 0)
		m.Set(2, 3, 0)
		m.Set(3, 3, 0)
	}
	return m
}

func (m *MatRxC) Compose(translate, rotate, scale []float32) Matrice {
	m.Translate(VecUnp(translate...))
	m.Rotate(QuatUnp(rotate...))
	m.Scale(VecUnp(scale...))
	return m
}

func (m *MatRxC) Decompose(translate, scale Vector, rotate Quaternion) Matrice {
	v := Vec3(0, 0, 0)

	v.Update(m.Get(0, 0), m.Get(0, 1), m.Get(1, 2))
	sx := v.Len()
	v.Update(m.Get(1, 0), m.Get(1, 1), m.Get(1, 2))
	sy := v.Len()
	v.Update(m.Get(2, 0), m.Get(2, 1), m.Get(2, 2))
	sz := v.Len()

	det := m.Determinant()
	if det < 0 {
		sx = -sx
	}

	if translate != nil {
		translate.Set(0, m.Get(0, 3)) //m[12]
		translate.Set(1, m.Get(1, 3)) //m[13]
		translate.Set(2, m.Get(2, 3)) //m[14]
	}

	invSX := 1 / sx
	invSY := 1 / sy
	invSZ := 1 / sz

	raw := m.Raw()
	raw[0] *= invSX
	raw[1] *= invSX
	raw[2] *= invSX

	raw[4] *= invSY
	raw[5] *= invSY
	raw[6] *= invSY

	raw[8] *= invSZ
	raw[9] *= invSZ
	raw[10] *= invSZ

	m.Update(raw...)

	if rotate != nil {
		//quaternion.SetFromRotationMatrix(&matrix)
	}

	if scale != nil {
		scale.Set(0, sx) //scale.X = sx
		scale.Set(1, sy) //scale.Y = sy
		scale.Set(2, sz) //scale.Z = sz
	}
	return m
}

func (m *MatRxC) LookAt(eye, target, up Vector) Matrice {
	v := Vec3(0, 0, 0)
	f, s, u := v.Clone(), v.Clone(), v.Clone()
	//f.SubVectors(target, eye).Normalize()
	//s.CrossVectors(&f, up).Normalize()
	//u.CrossVectors(&s, &f)

	raw := []float32{
		s.Get(0), u.Get(0), -(f.Get(0)), 0.0,
		s.Get(1), u.Get(1), -(f.Get(1)), 0.0,
		s.Get(2), u.Get(2), -(f.Get(2)), 0.0,
		-s.Dot(eye), -u.Dot(eye), f.Dot(eye), 1.0,
	}
	m.Update(raw...)
	return m
}

func (m *MatRxC) Frustrum(left, right, bottom, top, near, far float32) Matrice {
	f := []float32{
		2 * near / (right - left),
		0,
		0,
		0,
		0,
		2 * near / (top - bottom),
		0,
		0,
		(right + left) / (right - left),
		(top + bottom) / (top - bottom),
		-(far + near) / (far - near),
		-1,
		0,
		0,
		-(2 * far * near) / (far - near),
		0,
	}
	m.tag = MAT4
	m.v = f
	return m
}

func (m *MatRxC) Perspective(fov, aspect, near, far float32) Matrice {
	ymax := near * Tan(DegToRad(fov*0.5))
	ymin := -ymax
	xmin := ymin * aspect
	xmax := ymax * aspect
	return m.Frustrum(xmin, xmax, ymin, ymax, near, far)
}

func (m *MatRxC) Orthographic(left, right, top, bottom, near, far float32) Matrice {
	w := right - left
	h := top - bottom
	p := far - near

	x := (right + left) / w
	y := (top + bottom) / h
	z := (far + near) / p

	f := []float32{
		2 / w, 0, 0, 0,
		0, 2 / h, 0, 0,
		0, 0, -2 / p, 0,
		-x, -y, -z, 1,
	}
	m.tag = MAT4
	m.v = f
	return m
}

func ExpectedMatrice(k string, v interface{}) (Matrice, bool) {
	if mat, ok := v.(Matrice); ok {
		switch mat.Tag() {
		case MAT2:
			if k == MAT2 {
				return mat, true
			}
		case MAT3:
			if k == MAT3 {
				return mat, true
			}
		case MAT4:
			if k == MAT4 {
				return mat, true
			}
		}
	}
	return nil, false
}

func ToMatrice(L *l.LState, v l.LValue) Matrice {
	return toMatrice(L, v, defaultToErrMsgs("Matrice"))
}

func toMatrice(L *l.LState, v l.LValue, msg toErrMsgs) Matrice {
	if vv, ok := v.(*l.LUserData); ok {
		val := vv.Value
		if mat, ok := val.(Matrice); ok {
			return mat
		}
		L.RaiseError(msg[0], vv)
	}
	L.RaiseError(msg[1], v)
	return nil
}

func diagN(dst *MatRxC, diag *VecN) *MatRxC {
	ln := len(diag.v)
	r, c := ln, ln
	tag := matKindFromDimensions(r, c)
	dst = dst.reshape(tag, r, c)

	n := ln
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				dst.Set(i, j, diag.v[i])
			} else {
				dst.Set(i, j, 0)
			}
		}
	}

	return dst
}

func copyMatMN(dst, src *MatRxC) {
	if dst == nil || src == nil {
		return
	}
	dst.reshape(src.tag, src.r, src.c)
	copy(dst.v, src.v)
}

func flattenToArrayOffset(m *MatRxC, array []float32, offset int) []float32 {
	copy(array[offset:], m.v[:])
	return array
}

func IdentityMatrix(k string) Matrice {
	var m *MatRxC
	var d []float32
	switch k {
	case MAT2:
		m = Mat2()
		d = diagonals(1, 2)
	case MAT3:
		m = Mat3()
		d = diagonals(1, 3)
	case MAT4:
		m = Mat4()
		d = diagonals(1, 4)
	}
	return diagN(m, NewVecNFrom("", d))
}

func Mat2(v ...float32) *MatRxC {
	return newMatrixFromData(MAT2, v, 2, 2)
}

func Mat3(v ...float32) *MatRxC {
	return newMatrixFromData(MAT3, v, 3, 3)
}

func Mat4(v ...float32) *MatRxC {
	return newMatrixFromData(MAT4, v, 4, 4)
}

func lMat(k string) l.LGFunction {
	return func(L *l.LState) int {
		m := UnpackToMat(L, 1, k, false)
		push := func(u *l.LUserData) {
			u.Value = m
		}
		lua.PushNewUserData(L, push, k)
		return 1
	}
}

func initMatrix(k string) *MatRxC {
	var row, column int
	switch k {
	case MAT2:
		row, column = 2, 2
	case MAT3:
		row, column = 3, 3
	case MAT4:
		row, column = 4, 4
	}
	return newMatrix(k, row, column)
}

func diagonals(v float32, l int) []float32 {
	var ret []float32
	for i := 1; i <= l; i++ {
		ret = append(ret, v)
	}
	return ret
}

func rowColumnIdx(r, c, l int) (int, int) {
	if r >= l {
		c = c + 1
		r = 0
	}
	return r, c
}

func raiseErrorIf(L *l.LState, k string, i, t, c, l int) {
	if i < t {
		L.RaiseError("Too many arguments to %s", k)
	}
	if c < l {
		L.RaiseError("%s constructor arguments have insufficient components", k)
	}
}

func UnpackToMat(L *l.LState, from int, k string, ignore bool) Matrice {
	limit := establishLimit(k)
	m := initMatrix(k)
	top := L.GetTop()
	switch {
	case top == 1:
		unp := L.Get(from)
		switch unp.(type) {
		case l.LNumber:
			num := float32(L.CheckNumber(1))
			var diag []float32
			switch k {
			case MAT2:
				diag = diagonals(num, 2)
				m = diagN(m, NewVecNFrom("", diag))
				goto RETURN
			case MAT3:
				diag = diagonals(num, 3)
				m = diagN(m, NewVecNFrom("", diag))
				goto RETURN
			case MAT4:
				diag = diagonals(num, 4)
				m = diagN(m, NewVecNFrom("", diag))
				goto RETURN
			}
		case *l.LUserData:
			ud := unp.(*l.LUserData)
			d := ud.Value
			switch dt := d.(type) {
			case *MatRxC:
				switch dt.Tag() {
				case MAT2:
					copyMatMN(m, dt)
					goto RETURN
				case MAT3:
					copyMatMN(m, dt)
					goto RETURN
				case MAT4:
					copyMatMN(m, dt)
					goto RETURN
				}
			}
		}
	default:
		row, column := 0, 0
		for i := from; i <= top; i++ {
			unp := L.Get(i)
			switch unp.(type) {
			case l.LNumber:
				m.Set(row, column, float32(L.CheckNumber(i)))
				row = row + 1
				row, column = rowColumnIdx(row, column, limit)
				if column >= limit {
					raiseErrorIf(L, k, i, top, column, limit)
					goto RETURN
				}
			case *l.LUserData:
				unpToMat := func(m *MatRxC, v Vector, vl int) bool {
					for ii := 1; ii <= vl; ii++ {
						m.Set(row, column, v.Get(ii-1))
						row = row + 1
						row, column = rowColumnIdx(row, column, limit)
						if column >= limit {
							raiseErrorIf(L, k, i, top, column, limit)
							return true
						}
						continue
					}
					return false
				}
				ud := unp.(*l.LUserData)
				d := ud.Value
				switch dt := d.(type) {
				case *VecN:
					switch dt.Tag() {
					case VEC2:
						if unpToMat(m, dt, 2) {
							goto RETURN
						}
					case VEC3:
						if unpToMat(m, dt, 3) {
							goto RETURN
						}
					case VEC4:
						if unpToMat(m, dt, 4) {
							goto RETURN
						}
					}
				}
			default:
				if !ignore {
					L.RaiseError("%v cannot be passed to matrix", unp)
				}
			}
		}
	}
RETURN:
	return m
}

func matIndex(t *lua.Table, _ string) l.LGFunction {
	k := t.Name
	limit := establishLimit(k)
	vecFunc := func(limit, column int, m Matrice) *VecN {
		var list []float32
		for i := 1; i <= limit; i++ {
			list = append(list, m.Get(i-1, column-1))
		}
		return NewVecNFrom("", list)
	}
	return func(L *l.LState) int {
		ud := L.CheckUserData(1)
		if m, ok := ExpectedMatrice(k, ud.Value); ok {
			col := L.CheckInt(2)
			vec := vecFunc(limit, col, m)
			switch k {
			case MAT2:
				push := vec
				pfn := func(u *l.LUserData) {
					u.Value = push
				}
				lua.PushNewUserData(L, pfn, VEC2)
			case MAT3:
				push := vec
				pfn := func(u *l.LUserData) {
					u.Value = push
				}
				lua.PushNewUserData(L, pfn, VEC3)
			case MAT4:
				push := vec
				pfn := func(u *l.LUserData) {
					u.Value = push
				}
				lua.PushNewUserData(L, pfn, VEC4)
			}
			return 1
		}
		L.RaiseError("%s expected", k)
		return 0
	}
}

func MatOfSize(sz int, m []float32) (Matrice, error) {
	var ret Matrice
	var err error
	var row, column int
	switch sz {
	case 4:
		rm := initMatrix(MAT2)
		for i := 0; i <= sz-1; i++ {
			rm.Set(row, column, m[i])
			row = row + 1
			row, column = rowColumnIdx(row, column, 2)
		}
		ret = rm
	case 9:
		rm := initMatrix(MAT3)
		for i := 0; i <= sz-1; i++ {
			rm.Set(row, column, m[i])
			row = row + 1
			row, column = rowColumnIdx(row, column, 3)
		}
		ret = rm
	case 16:
		rm := initMatrix(MAT4)
		for i := 0; i <= sz-1; i++ {
			rm.Set(row, column, m[i])
			row = row + 1
			row, column = rowColumnIdx(row, column, 4)
		}
		ret = rm
	default:
		err = NoXOfSizeError("matrice", sz)
	}
	return ret, err
}

func matIMF(L *l.LState, k string) innerMathFunc {
	return func(v1, v2 []float32, imr innerMathRun) int {
		size, dst := imr(v1, v2)
		nm, err := MatOfSize(size, dst)
		if err != nil {
			L.RaiseError("unable to perform matrice %s: %s", k, err.Error())
		}
		fn := func(u *l.LUserData) {
			u.Value = nm
		}
		mt := establishKindOfLimit("matrice", size)
		lua.PushNewUserData(L, fn, mt)
		return 1
	}
}

func matVaryParams(L *l.LState) ([]float32, []float32) {
	var rv1, rv2 []float32

	var mpos1ErrMsgs = toErrMsgs{
		"first position param %s is not UserData",
		"first position param %s is not Matrice",
	}

	v1 := toMatrice(L, L.Get(1), mpos1ErrMsgs)
	rv1 = v1.Raw()

	v2 := L.Get(2)
	switch v2.(type) {
	case l.LNumber:
		c := Pf32(L, 2)
		rv2 = []float32{c}
	case *l.LUserData:
		v2v := ToMatrice(L, v2)
		rv2 = v2v.Raw()
	default:
		L.RaiseError("%s not a matrice applicable math type", v2)
	}

	return rv1, rv2
}

func matInner(imr innerMathRun) inner {
	return func(L *l.LState, vm innerMathFunc) int {
		v1, v2 := matVaryParams(L)
		return vm(v1, v2, imr)
	}
}

func matAdd(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		imr := func(v1, v2 []float32) (int, []float32) {
			var size int
			var dst []float32
			switch len(v2) {
			case 1:
				size = len(v1)
				dst = make([]float32, size)
				for idx, item := range v1 {
					dst[idx] = item + v2[0]
				}
			default:
				size = intMin(len(v1), len(v2))
				dst = make([]float32, size)
				for i := 1; i <= size; i++ {
					idx := i - 1
					dst[idx] = v1[idx] + v2[idx]
				}
			}
			return size, dst
		}
		return runInner(
			L,
			matInner(imr),
			matIMF(L, "addition"),
		)
	}
}

func matSub(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		imr := func(v1, v2 []float32) (int, []float32) {
			var size int
			var dst []float32
			switch len(v2) {
			case 1:
				size = len(v1)
				dst = make([]float32, size)
				for idx, item := range v1 {
					dst[idx] = item - v2[0]
				}
			default:
				size = intMin(len(v1), len(v2))
				dst = make([]float32, size)
				for i := 1; i <= size; i++ {
					idx := i - 1
					dst[idx] = v1[idx] - v2[idx]
				}
			}
			return size, dst
		}
		return runInner(
			L,
			matInner(imr),
			matIMF(L, "subtraction"),
		)
	}
}

func matDimensions(rm []float32) (int, int) {
	var kind string
	var rows, cols int
	switch len(rm) {
	case 4:
		kind = MAT2
	case 9:
		kind = MAT3
	case 16:
		kind = MAT4
	}
	switch kind {
	case MAT2:
		rows, cols = 2, 2
	case MAT3:
		rows, cols = 3, 3
	case MAT4:
		rows, cols = 4, 4

	}
	return rows, cols
}

func matKindFromDimensions(r, c int) string {
	var ret string
	switch r * c {
	case 4:
		ret = MAT2
	case 9:
		ret = MAT3
	case 16:
		ret = MAT4
	}
	return ret
}

func matMul(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		imr := func(v1, v2 []float32) (int, []float32) {
			var size int
			var dst []float32
			switch len(v2) {
			case 1:
				size = len(v1)
				dst = make([]float32, size)
				for idx, item := range v1 {
					dst[idx] = item * v2[0]
				}
			default:
				m1r, m1c := matDimensions(v1)
				m2r, m2c := matDimensions(v2)
				dstKind := matKindFromDimensions(m1r, m2c)
				switch dstKind {
				case MAT2:
					size = 4
				case MAT3:
					size = 9
				case MAT4:
					size = 16
				}
				dst = make([]float32, size)
				if m1r != m2c {
					L.RaiseError("matrices cannot be multiplied: #first rows != #second columns")

				}
				for r1 := 0; r1 < m1r; r1++ {
					for c2 := 0; c2 < m2c; c2++ {
						dst[c2*m1r+r1] = 0
						for i := 0; i < m1c; i++ {
							dst[c2*m1r+r1] += v1[i*m1r+r1] * v2[c2*m2r+i]
						}
					}
				}
			}
			return size, dst
		}
		return runInner(
			L,
			matInner(imr),
			matIMF(L, "multiplication"),
		)
	}
}

func matDiv(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		//imr := func(v1, v2 []float32) (int, []float32) {
		//	size := 0 //intMin(len(v1), len(v2))
		//	dst := make([]float32, size)
		//	//for i := 1; i <= size; i++ {
		//	//	idx := i - 1
		//	//	dst[idx] = v1[idx] + v2[idx]
		//	//}
		//	return size, dst
		//}
		//return runInner(
		//	L,
		//	matInner(imr),
		//	matIMF(L, "division"),
		//)
		L.Push(l.LString("UNIMPLEMENTED"))
		return 1
	}
}

func matUnm(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		//k := t.Name
		//ud := L.CheckUserData(1)
		//if v, ok := ExpectedVector(k, ud.Value); ok {
		//	lt := establishLimit(k)
		//	for i := 0; i <= lt-1; i++ {
		//		nv := v.Get(i)
		//		nv = nv - 1
		//		v.Set(i, nv)
		//	}
		//	L.Push(ud)
		//	return 1
		//}
		//L.RaiseError("__unm error, %s expected", k)
		L.Push(l.LString("UNIMPLEMENTED"))
		return 1
	}
}

func matLen(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		k := t.Name
		ud := L.CheckUserData(1)
		if m, ok := ExpectedMatrice(k, ud.Value); ok {
			cl := m.Cols()
			L.Push(l.LNumber(cl))
			return 1
		}
		L.RaiseError("__len error, %s expected", k)
		return 0
	}
}

//return raw list, return item by col row position
func matCall(t *lua.Table, _ string) l.LGFunction {
	return func(L *l.LState) int {
		//spew.Dump(L.GetTop())
		return 0
	}
}

var matMeta = []*lua.LMetaFunc{
	{"__index", matIndex},
	lua.ImmutableNewIdx(),
	{"__add", matAdd},
	{"__sub", matSub},
	{"__mul", matMul},
	{"__div", matDiv},
	{"__unm", matUnm},
	{"__len", matLen},
	{"__call", matCall},
}
