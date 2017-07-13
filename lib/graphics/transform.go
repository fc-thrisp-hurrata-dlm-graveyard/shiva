package graphics

import "github.com/Laughs-In-Flowers/shiva/lib/math"

type Translate struct {
	math.Vector
}

func NewTranslate(v math.Vector) *Translate {
	return &Translate{
		v,
	}
}

func (t *Translate) Provide(p Provider) {}

type Scale struct {
	math.Vector
}

func NewScale(v math.Vector) *Scale {
	return &Scale{
		v,
	}
}

func (s *Scale) Provide(p Provider) {}

type Rotate struct {
	math.Quaternion
}

func NewRotate(q math.Quaternion) *Rotate {
	return &Rotate{
		q,
	}
}

func (r *Rotate) Provide(p Provider) {}
