package math

import (
	"math"

	"github.com/Laughs-In-Flowers/shiva/lib/lua"

	l "github.com/yuin/gopher-lua"
)

type MathType interface {
	Tag() string
}

type Manipulator interface {
	Raw() []float32
	Update(...float32)
}

type Getter interface {
	Get(int) float32
	GetStr(string) float32
}

type MultiGetter interface {
	Get(int, int) float32
	GetStr(string, string) float32
}

type Setter interface {
	Set(int, float32)
	SetStr(string, float32)
}

type MultiSetter interface {
	Set(int, int, float32)
	SetStr(string, string, float32)
}

type Length interface {
	RawLen() int
	Len() float32
}

const Pi = math.Pi
const degreeToRadiansFactor = math.Pi / 180
const radianToDegreesFactor = 180.0 / math.Pi

var Infinity = float32(math.Inf(1))

func DegToRad(degrees float32) float32 {
	return degrees * degreeToRadiansFactor
}

func RadToDeg(radians float32) float32 {
	return radians * radianToDegreesFactor
}

// Clamp value to range <a, b>
func Clamp(x, a, b float32) float32 {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

func ClampInt(x, a, b int) int {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

// Clamp value to range <a, inf)
func ClampBotton(x, a float32) float32 {
	if x < a {
		return a
	}
	return x
}

func Abs(v float32) float32 {
	return float32(math.Abs(float64(v)))
}

func Acos(v float32) float32 {
	return float32(math.Acos(float64(v)))
}

func Asin(v float32) float32 {
	return float32(math.Asin(float64(v)))
}

func Atan(v float32) float32 {
	return float32(math.Atan(float64(v)))
}

func Atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func Ceil(v float32) float32 {
	return float32(math.Ceil(float64(v)))
}

func Cos(v float32) float32 {
	return float32(math.Cos(float64(v)))
}

func Floor(v float32) float32 {
	return float32(math.Floor(float64(v)))
}

func Inf(sign int) float32 {
	return float32(math.Inf(sign))
}

func Round(v float32) float32 {
	return Floor(v + 0.5)
}

func IsNaN(v float32) bool {
	return math.IsNaN(float64(v))
}

func Sin(v float32) float32 {
	return float32(math.Sin(float64(v)))
}

func Sqrt(v float32) float32 {
	return float32(math.Sqrt(float64(v)))
}

func Max(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
}

func Min(a, b float32) float32 {
	return float32(math.Min(float64(a), float64(b)))
}

func Mod(a, b float32) float32 {
	return float32(math.Mod(float64(a), float64(b)))
}

func NaN() float32 {
	return float32(math.NaN())
}

func Pow(a, b float32) float32 {
	return float32(math.Pow(float64(a), float64(b)))
}

func Tan(v float32) float32 {
	return float32(math.Tan(float64(v)))
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
