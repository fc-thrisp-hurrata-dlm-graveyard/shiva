package opengl45

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/xrror"
)

func New(debug bool) graphics.Provider {
	if debug {
		return &OGL45DEBUG{"opengl", 4, 5, 0, debug}
	}
	return &OGL45{"opengl", 4, 5, 0, debug}
}

var InitFailure = xrror.Xrror("Failed to initialize OpenGL 4.5(debug %d): %s").Out

func init() {
	graphics.Register("opengl4.5", New)
}
