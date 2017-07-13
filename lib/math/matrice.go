package math

import (
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/go-gl/mathgl/mgl32"

	l "github.com/yuin/gopher-lua"
)

type Matrice interface {
	MathType
	Raw() []float32
	At(int, int) float32
	Set(int, int, float32)
	NumCols() int
}

func ExpectedMatrice(k string, v interface{}) (Matrice, bool) {
	if mat, ok := v.(Matrice); ok {
		switch mat.(type) {
		case *mat2:
			if k == MAT2 {
				return mat, true
			}
		case *mat3:
			if k == MAT3 {
				return mat, true
			}
		case *mat4:
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

type mat2 struct {
	*mgl32.MatMxN
}

func (*mat2) Tag() string {
	return MAT2
}

type mat3 struct {
	*mgl32.MatMxN
}

func (*mat3) Tag() string {
	return MAT3
}

type mat4 struct {
	*mgl32.MatMxN
}

func (*mat4) Tag() string {
	return MAT4
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

func initMatrix(k string) *mgl32.MatMxN {
	var ret *mgl32.MatMxN
	switch k {
	case MAT2:
		ret = mgl32.NewMatrix(2, 2)
	case MAT3:
		ret = mgl32.NewMatrix(3, 3)
	case MAT4:
		ret = mgl32.NewMatrix(4, 4)
	}
	return ret
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
				m = mgl32.DiagN(m, mgl32.NewVecNFromData(diag))
				goto RETURN
			case MAT3:
				diag = diagonals(num, 3)
				m = mgl32.DiagN(m, mgl32.NewVecNFromData(diag))
				goto RETURN
			case MAT4:
				diag = diagonals(num, 4)
				m = mgl32.DiagN(m, mgl32.NewVecNFromData(diag))
				goto RETURN
			}
		case *l.LUserData:
			ud := unp.(*l.LUserData)
			d := ud.Value
			switch dt := d.(type) {
			case *mat2:
				mgl32.CopyMatMN(m, dt.MatMxN)
				goto RETURN
			case *mat3:
				mgl32.CopyMatMN(m, dt.MatMxN)
				goto RETURN
			case *mat4:
				mgl32.CopyMatMN(m, dt.MatMxN)
				goto RETURN
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
				unpToMat := func(m *mgl32.MatMxN, v Vector, vl int) bool {
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
				case *vec2:
					if unpToMat(m, dt, 2) {
						goto RETURN
					}
				case *vec3:
					if unpToMat(m, dt, 3) {
						goto RETURN
					}
				case *vec4:
					if unpToMat(m, dt, 4) {
						goto RETURN
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
	var ret Matrice
	switch k {
	case MAT2:
		ret = &mat2{m}
	case MAT3:
		ret = &mat3{m}
	case MAT4:
		ret = &mat4{m}
	}
	return ret
}

func matIndex(t *lua.Table, _ string) l.LGFunction {
	k := t.Name
	limit := establishLimit(k)
	vecFunc := func(limit, column int, m Matrice) *VecN {
		var list []float32
		for i := 1; i <= limit; i++ {
			list = append(list, m.At(i-1, column-1))
		}
		return NewVecNFrom(list)
	}
	return func(L *l.LState) int {
		ud := L.CheckUserData(1)
		if m, ok := ExpectedMatrice(k, ud.Value); ok {
			col := L.CheckInt(2)
			vec := vecFunc(limit, col, m)
			switch k {
			case MAT2:
				push := &vec2{vec}
				pfn := func(u *l.LUserData) {
					u.Value = push
				}
				lua.PushNewUserData(L, pfn, VEC2)
			case MAT3:
				push := &vec3{vec}
				pfn := func(u *l.LUserData) {
					u.Value = push
				}
				lua.PushNewUserData(L, pfn, VEC3)
			case MAT4:
				push := &vec4{vec}
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
		ret = &mat2{rm}
	case 9:
		rm := initMatrix(MAT3)
		for i := 0; i <= sz-1; i++ {
			rm.Set(row, column, m[i])
			row = row + 1
			row, column = rowColumnIdx(row, column, 3)
		}
		ret = &mat3{rm}
	case 16:
		rm := initMatrix(MAT4)
		for i := 0; i <= sz-1; i++ {
			rm.Set(row, column, m[i])
			row = row + 1
			row, column = rowColumnIdx(row, column, 4)
		}
		ret = &mat4{rm}
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
			cl := m.NumCols()
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
