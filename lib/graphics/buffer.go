package graphics

import (
	"unsafe"

	"github.com/Laughs-In-Flowers/shiva/lib/math"
)

type BuffAttrib struct {
	Name string
	Size int32
}

type Buff struct {
	p      Provider
	handle Buffer
	usage  Enum
	update bool
	buffer math.AF32
	a      []BuffAttrib
}

func NewBuff() *Buff {
	return &Buff{
		usage:  STATIC_DRAW,
		update: true,
		a:      make([]BuffAttrib, 0),
	}
}

func (b *Buff) Attrib(name string) *BuffAttrib {
	for _, a := range b.a {
		if name == a.Name {
			return &a
		}
	}
	return nil
}

func (b *Buff) AttribAt(idx int) *BuffAttrib {
	return &b.a[idx]
}

func (b *Buff) AddAttrib(name string, size int32) *Buff {
	b.a = append(b.a, BuffAttrib{name, size})
	return b
}

func (b *Buff) AttribCount() int {
	return len(b.a)
}

func (b *Buff) SetBuffer(r []float32) *Buff {
	b.buffer = r
	b.update = true
	return b
}

func (b *Buff) SetUsage(u Enum) {
	b.usage = u
}

func (b *Buff) Buffer() *math.AF32 {
	return &b.buffer
}

func (b *Buff) Update() {
	b.update = true
}

func (b *Buff) Stride() int {
	stride := 0
	elsize := int(unsafe.Sizeof(float32(0)))
	for _, attrib := range b.a {
		stride += elsize * int(attrib.Size)
	}
	return stride
}

//var AttribLocError = xrror.Xrror("unable to retrieve attribute %s location from program %v").Out

func (b *Buff) Provide(p Provider) {
	bb := b.buffer

	if bb.Bytes() == 0 {
		return
	}

	if b.p == nil {
		b.handle = p.GenBuffer()
		p.BindBuffer(ARRAY_BUFFER, b.handle)
		stride := b.Stride()
		var items uint32 = 0
		var offset uint32 = 0
		elsize := int32(unsafe.Sizeof(float32(0)))
		for _, attrib := range b.a {
			var prog Program
			prog = p.CurrentProgram()
			if prog == 0 {
			}
			loc := p.GetAttribLocation(prog, attrib.Name)
			if loc < 0 {
				continue //return AttribLocError(attrib.Name, prog)
			}
			p.EnableVertexAttribArray(uint32(loc))
			p.VertexAttribPointer(uint32(loc), attrib.Size, FLOAT, false, int32(stride), p.Ptr(&offset))
			items += uint32(attrib.Size)
			offset = uint32(elsize) * items
		}

		b.p = p
	}

	if !b.update {
		return
	}

	p.BindBuffer(ARRAY_BUFFER, b.handle)
	p.BufferData(ARRAY_BUFFER, bb.Bytes(), p.Ptr(&bb[0]), b.usage)
	b.update = false
}

func (b *Buff) Close() {
	if p := b.p; p != nil {
		p.DeleteBuffer(b.handle)
	}
	b.p = nil
}

/*
type AttrType int

// List of all possible attribute types.
const (
	Int AttrType = iota
	Float
	Vec2
	Vec3
	Vec4
	Mat2
	Mat23
	Mat24
	Mat3
	Mat32
	Mat34
	Mat4
	Mat42
	Mat43
)

// Size returns the size of a type in bytes.
func (at AttrType) Size() int {
	switch at {
	case Int:
		return 4
	case Float:
		return 4
	case Vec2:
		return 2 * 4
	case Vec3:
		return 3 * 4
	case Vec4:
		return 4 * 4
	case Mat2:
		return 2 * 2 * 4
	case Mat23:
		return 2 * 3 * 4
	case Mat24:
		return 2 * 4 * 4
	case Mat3:
		return 3 * 3 * 4
	case Mat32:
		return 3 * 2 * 4
	case Mat34:
		return 3 * 4 * 4
	case Mat4:
		return 4 * 4 * 4
	case Mat42:
		return 4 * 2 * 4
	case Mat43:
		return 4 * 3 * 4
	default:
		panic("size of vertex attribute type: invalid type")
	}
}
*/
