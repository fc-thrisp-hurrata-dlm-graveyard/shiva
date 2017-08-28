package graphics

import "fmt"

type Uniform interface {
	Update(...float32)
	Transfer(Provider)
	TransferIdx(Provider, int)
}

type innerTransfer func(Provider, int32, []float32)

type uniform struct {
	key    string
	keyIdx string
	idx    int
	value  []float32
	trf    innerTransfer
}

func newUniform(key string, length int, trf innerTransfer) *uniform {
	return &uniform{
		key:   key,
		value: make([]float32, length),
		trf:   trf,
	}
}

func (u *uniform) Update(v ...float32) {
	u.value = v
}

func (u *uniform) location(p Provider) int32 {
	return p.GetUniformCurrentLocation(u.key)
}

func (u *uniform) locationIdx(p Provider, idx int) int32 {
	if u.keyIdx == "" || u.idx != idx {
		u.keyIdx = fmt.Sprintf("%s[%d]", u.key, idx)
		u.idx = idx
	}
	return p.GetUniformCurrentLocation(u.keyIdx)
}

func (u *uniform) Transfer(p Provider) {
	loc := u.location(p)
	u.trf(p, loc, u.value)
}

func (u *uniform) TransferIdx(p Provider, idx int) {
	loc := u.locationIdx(p, idx)
	u.trf(p, loc, u.value)
}

func Uniform1i(key string) Uniform {
	return newUniform(key, 1, func(p Provider, loc int32, v []float32) {
		p.Uniform1i(loc, int32(v[0]))
	})
}

func Uniform1f(key string) Uniform {
	return newUniform(key, 1, func(p Provider, loc int32, v []float32) {
		p.Uniform1f(loc, v[0])
	})
}

func Uniform2f(key string) Uniform {
	return newUniform(key, 2, func(p Provider, loc int32, v []float32) {
		//p.Uniform2f(u.Location(p), uni.v0)
	})
}

func Uniform3f(key string) Uniform {
	return newUniform(key, 3, func(p Provider, loc int32, v []float32) {
		//p.Uniform3f(u.Location(p), uni.v0)
	})
}

func Uniform4f(key string) Uniform {
	return newUniform(key, 4, func(p Provider, loc int32, v []float32) {
		//p.Uniform4f(u.Location(p), uni.v0)
	})
}

func UniformMatrix3fv(key string) Uniform {
	return newUniform(key, 9, func(p Provider, loc int32, v []float32) {
		//p.UniformMatrix3fv(u.Location(p), uni.v0)
	})
}

func UniformMatrix4fv(key string) Uniform {
	return newUniform(key, 16, func(p Provider, loc int32, v []float32) {
		p.UniformMatrix4fv(loc, 1, false, v)
	})
}

func Uniform4fv(key string) Uniform {
	return newUniform(key, 16, func(p Provider, loc int32, v []float32) {
		//p.Uniform4f(u.Location(p), uni.v0)
	})
}
