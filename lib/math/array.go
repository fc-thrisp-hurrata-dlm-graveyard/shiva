package math

import "unsafe"

// array of float32
type AF32 []float32

func NewAF32(size, capacity int) AF32 {
	return make([]float32, size, capacity)
}

// Bytes returns the size of the array in bytes
func (a *AF32) Bytes() int {
	return len(*a) * int(unsafe.Sizeof(float32(0)))
}

// Size returns the number of float32 elements in the array
func (a *AF32) Size() int {
	return len(*a)
}

// Len returns the number of float32 elements in the array
func (a *AF32) Len() int {
	return len(*a)
}

func (a *AF32) Append(v ...float32) {
	*a = append(*a, v...)
}

func (a AF32) Set(pos int, v ...float32) {
	for i := 0; i < len(v); i++ {
		a[pos+i] = v[i]
	}
}

// array of uint32
type AU32 []uint32

func NewAU32(size, capacity int) AU32 {
	return make([]uint32, size, capacity)
}

// Bytes returns the size of the array in bytes
func (a *AU32) Bytes() int {
	return len(*a) * int(unsafe.Sizeof(float32(0)))
}

// Size returns the number of uint32 elements in the array
func (a *AU32) Size() int {
	return len(*a)
}

// Len returns the number of uint32 elements in the array
func (a *AU32) Len() int {
	return len(*a)
}

func (a *AU32) Append(v ...uint32) {
	*a = append(*a, v...)
}

func (a AU32) Set(pos int, v ...uint32) {
	for i := 0; i < len(v); i++ {
		a[pos+i] = v[i]
	}
}
