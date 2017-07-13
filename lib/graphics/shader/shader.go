package shader

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/xrror"
)

var noShaderError = xrror.Xrror("%s not mappable to a graphics shader type").Out

func toGLShaderEnum(key string) (graphics.Enum, error) {
	switch strings.ToLower(key) {
	case "vertex":
		return graphics.VERTEX_SHADER, nil
	case "geometry":
		return graphics.GEOMETRY_SHADER, nil
	case "fragment":
		return graphics.FRAGMENT_SHADER, nil
	}
	return 0, noShaderError(key)
}

func NewShader(p graphics.Provider, pr *Profile, sr Shaderer, kind, tag string) (graphics.Shader, error) {
	raw, err := renderShaderTemplate(sr, pr, tag)
	if err != nil {
		return 0, err
	}
	return buildShader(p, kind, raw)
}

func renderShaderTemplate(sr Shaderer, pr *Profile, tag string) (string, error) {
	b := new(bytes.Buffer)
	err := sr.Render(b, tag, pr)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func buildShader(p graphics.Provider, kind, source string) (graphics.Shader, error) {
	var sge graphics.Enum
	var err error
	sge, err = toGLShaderEnum(kind)
	if err != nil {
		return 0, err
	}
	var handle graphics.Shader
	handle, err = compileShader(p, source, sge)
	if err != nil {
		return 0, err
	}
	return handle, nil
}

func compileShader(p graphics.Provider, raw string, shaderType graphics.Enum) (graphics.Shader, error) {
	shader := p.CreateShader(shaderType)
	p.ShaderSource(shader, raw)

	p.CompileShader(shader)

	var status int32
	p.GetShaderiv(shader, graphics.COMPILE_STATUS, &status)
	if status == graphics.FALSE {
		msg := p.GetShaderInfoLog(shader)
		errmsg := fmt.Errorf("failed to compile %v: %v", raw, msg)
		return 0, errmsg
	}
	return shader, nil
}
