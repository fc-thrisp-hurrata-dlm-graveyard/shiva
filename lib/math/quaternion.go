package math

import (
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/go-gl/mathgl/mgl32"

	l "github.com/yuin/gopher-lua"
)

type Quaternion interface {
	W() float32
	X() float32
	Y() float32
	Z() float32
}

func expectedQuat(v interface{}) (Quaternion, bool) {
	return nil, false
}

type quat struct {
	q *mgl32.Quat
}

func (*quat) Tag() string {
	return QUAT
}

func (q *quat) W() float32 {
	return q.q.W
}

func Quat(w float32, v [3]float32) *quat {
	return &quat{
		&mgl32.Quat{w, mgl32.Vec3(v)},
	}
}

var defaultQVec = [3]float32{0, 0, 1}

func lQuat(L *l.LState) int {
	var q *quat
	top := L.GetTop()
	switch top {
	case 1:
		v := L.Get(1)
		switch v.(type) {
		case l.LNumber:
			num := Pf32(L, 1)
			q = Quat(num, defaultQVec)
		case *l.LUserData:
			// Mat3 or Mat4
		}
	case 2:
		var pos1 float32
		var pos2 [3]float32
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

		q = Quat(pos1, pos2)
	case 3:
		pos2, pos3, pos4 := Pf32(L, 1), Pf32(L, 2), Pf32(L, 3)
		q = Quat(0, [3]float32{pos2, pos3, pos4})
	case 4:
		pos1, pos2, pos3, pos4 := Pf32(L, 1), Pf32(L, 2), Pf32(L, 3), Pf32(L, 4)
		q = Quat(pos1, [3]float32{pos2, pos3, pos4})
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
		if v, ok := expectedQuat(ud.Value); ok {
			req := L.Get(2)
			switch req.String() {
			case "x", "pitch":
				L.Push(l.LNumber(v.X()))
				return 1
			case "y", "roll":
				L.Push(l.LNumber(v.Y()))
				return 1
			case "z", "yaw":
				L.Push(l.LNumber(v.Z()))
				return 1
			case "w", "angle":
				L.Push(l.LNumber(v.W()))
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
