package graphics

import "fmt"

type Uniform interface {
	Get() []float32
	Set(...float32)
	Transfer(Provider)
	TransferIdx(Provider, int)
	Location(Provider) int32
	LocationIdx(Provider, int) int32
}

type innerTransfer func(Provider, int32, *uniform)

type uniform struct {
	name    string
	nameIdx string
	idx     int
	v       []float32
	trf     innerTransfer
}

func (u *uniform) Get() []float32 {
	return u.v
}

func (u *uniform) Geti() int32 {
	return int32(u.v[0])
}

func (u *uniform) Getf() float32 {
	return u.v[0]
}

func (u *uniform) Set(v ...float32) {
	u.v = v
}

func (u *uniform) Location(p Provider) int32 {
	return p.GetCurrentUniformLocation(u.name)
}

func (u *uniform) LocationIdx(p Provider, idx int) int32 {
	if u.nameIdx == "" || u.idx != idx {
		u.nameIdx = fmt.Sprintf("%s[%d]", u.name, idx)
		u.idx = idx
	}
	return p.GetCurrentUniformLocation(u.nameIdx)
}

func (u *uniform) Transfer(p Provider) {
	loc := u.Location(p)
	u.trf(p, loc, u)
}

func (u *uniform) TransferIdx(p Provider, idx int) {
	loc := u.LocationIdx(p, idx)
	u.trf(p, loc, u)
}

func Uniform1i(name string) Uniform {
	return &uniform{
		name: name,
		trf: func(p Provider, loc int32, u *uniform) {
			p.Uniform1i(loc, u.Geti())
		},
		v: []float32{0},
	}
}

func Uniform1f(name string) Uniform {
	return &uniform{
		name: name,
		trf: func(p Provider, loc int32, u *uniform) {
			p.Uniform1f(loc, u.Getf())
		},
		v: []float32{0},
	}
}

func Uniform2f(name string) Uniform {
	return &uniform{
		name: name,
		trf: func(p Provider, loc int32, u *uniform) {
			//p.Uniform2f(u.Location(p), uni.v0)
		},
		v: []float32{0, 0},
	}
}

func Uniform3f(name string) Uniform {
	return &uniform{
		name: name,
		trf: func(p Provider, loc int32, u *uniform) {
			//p.Uniform3f(u.Location(p), uni.v0)
		},
		v: []float32{0, 0, 0},
	}
}

func Uniform4f(name string) Uniform {
	return &uniform{
		name: name,
		trf: func(p Provider, loc int32, u *uniform) {
			//p.Uniform4f(u.Location(p), uni.v0)
		},
		v: []float32{0, 0, 0, 0},
	}
}

func UniformMatrix3fv(name string) Uniform {
	return &uniform{
		name: name,
		trf: func(p Provider, loc int32, u *uniform) {
			//p.UniformMatrix3fv(u.Location(p), uni.v0)
		},
		v: []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func UniformMatrix4fv(name string) Uniform {
	return &uniform{
		name: name,
		trf: func(p Provider, loc int32, u *uniform) {
			p.UniformMatrix4fv(u.Location(p), 1, false, u.v)
		},
		v: []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func Uniform4fv(name string) Uniform {
	return &uniform{
		name: name,
		trf: func(p Provider, loc int32, u *uniform) {
			//p.Uniform4f(u.Location(p), uni.v0)
		},
	}
}
