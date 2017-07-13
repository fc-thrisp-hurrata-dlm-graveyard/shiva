package shader

import (
	"fmt"

	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/material"
)

type Prog struct {
	Tag, Version, F, G, V string
}

var defaultVersion string = "410 core"

var defaultProg = []*Prog{
	{"basic", defaultVersion, "fbasic", "", "vbasic"},
	{"standard", defaultVersion, "fstandard", "", "vstandard"},
}

type Profile struct {
	*Prog
	Independent bool
	UseLights   material.UseLights
	AmbientLightsMax, DirectionalLightsMax,
	PointLightsMax, SpotLightsMax, MaterialTexturesMax int
}

func (p *Profile) Equals(o *Profile) bool {
	switch {
	case p.Tag != o.Tag:
		return false
	case p.F != o.F && p.G != o.G && p.V != o.V:
		return false
	case o.Independent:
		return true
	case p.AmbientLightsMax == o.AmbientLightsMax &&
		p.DirectionalLightsMax == o.DirectionalLightsMax &&
		p.PointLightsMax == o.PointLightsMax &&
		p.SpotLightsMax == o.SpotLightsMax &&
		p.MaterialTexturesMax == o.MaterialTexturesMax:
		return true
	}
	return false
}

type Program struct {
	*Profile
	handle graphics.Program
}

func NewProgram(p graphics.Provider, pr *Profile, sr Shaderer, s ...graphics.Shader) (*Program, error) {
	var ss []graphics.Shader = s
	if len(ss) < 1 {
		var frag, geo, vert graphics.Shader
		var err error

		if pr.F != "" {
			frag, err = NewShader(p, pr, sr, "fragment", pr.F)
			if err != nil {
				return nil, err
			}
			ss = append(ss, frag)
		}

		if pr.G != "" {
			geo, err = NewShader(p, pr, sr, "geometry", pr.G)
			if err != nil {
				return nil, err
			}
			ss = append(ss, geo)
		}

		if pr.V != "" {
			vert, err = NewShader(p, pr, sr, "vertex", pr.V)
			if err != nil {
				return nil, err
			}
			ss = append(ss, vert)
		}
	}
	handle, err := compileProgram(p, ss...)
	if err != nil {
		return nil, err
	}

	return &Program{
		pr, handle,
	}, nil
}

func compileProgram(p graphics.Provider, shaders ...graphics.Shader) (graphics.Program, error) {
	programID := p.CreateProgram()

	for _, s := range shaders {
		p.AttachShader(programID, s)
	}

	p.LinkProgram(programID)

	var status int32
	p.GetProgramiv(programID, graphics.LINK_STATUS, &status)
	if status == graphics.FALSE {
		log := p.GetProgramInfoLog(programID)
		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	return programID, nil
}

func (p *Program) Provide(g graphics.Provider) {
	g.UseProgram(p.handle)
}

type Bind struct {
	Uniforms []graphics.Uniform
}

func NewBind(uu ...graphics.Uniform) *Bind {
	return &Bind{
		uu,
	}
}

func (b *Bind) Bind(p graphics.Provider, pr *Program) {
	//program := pr.handle
	//for _, u := range b.Uniforms {
	//u.Transfer(p, program)
	// OR TransferIdx ?
	//}
}
