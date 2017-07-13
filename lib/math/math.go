package math

import (
	"github.com/Laughs-In-Flowers/shiva/lib/lua"

	l "github.com/yuin/gopher-lua"
)

type MathType interface {
	Tag() string
}

type innerMathRun func([]float32, []float32) (int, []float32)

type innerMathFunc func([]float32, []float32, innerMathRun) int

type inner func(*l.LState, innerMathFunc) int

func runInner(L *l.LState, ii inner, im innerMathFunc) int {
	return ii(L, im)
}

const (
	VEC2 = "VEC2"
	VEC3 = "VEC3"
	VEC4 = "VEC4"
	MAT2 = "MAT2"
	MAT3 = "MAT3"
	MAT4 = "MAT4"
	QUAT = "QUAT"
)

var (
	mathMTs = []*lua.Table{
		lua.NewTable(VEC2, nil, vecMeta, nil, nil),
		lua.NewTable(VEC3, nil, vecMeta, nil, nil),
		lua.NewTable(VEC4, nil, vecMeta, nil, nil),
		lua.NewTable(MAT2, nil, matMeta, nil, nil),
		lua.NewTable(MAT3, nil, matMeta, nil, nil),
		lua.NewTable(MAT4, nil, matMeta, nil, nil),
		lua.NewTable(QUAT, nil, qutMeta, nil, nil),
	}

	expandMathFuncs = map[string]l.LGFunction{
		"vec2": lVec(VEC2),
		"vec3": lVec(VEC3),
		"vec4": lVec(VEC4),
		"mat2": lMat(MAT2),
		"mat3": lMat(MAT3),
		"mat4": lMat(MAT4),
		"quat": lQuat,
	}
)

func Module(L *lua.Lua) lua.Module {
	mopen := func(LL *l.LState) int {
		l.OpenMath(LL)
		L.RegisterFuncsOn("math", expandMathFuncs)
		// raise entire math module to use anywhere
		if mtb, err := L.GetModule("math"); err == nil {
			mtb.ForEach(func(k, v l.LValue) {
				key := k.String()
				LL.SetGlobal(key, v)
			})
		}
		for _, t := range mathMTs {
			mt := L.NewTypeMetatable(t.Name)
			for _, v := range t.Meta {
				key := v.Key
				fn := v.Value
				L.SetField(mt, key, L.NewClosure(fn(t, key)))
			}
		}
		return 1
	}
	return lua.NewModule("math", mopen)
}
