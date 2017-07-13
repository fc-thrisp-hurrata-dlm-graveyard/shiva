package texture

import (
	"image"

	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
)

type Texturer interface {
	graphics.Initializer
	graphics.Closer
	///graphics.Visiblizer
	graphics.RefCounter
}

type Texture interface {
	Texturer
	Renderer
}

type texture2D struct {
	p            graphics.Provider // Pointer to OpenGL state
	refCount     int               // Current number of references
	handle       graphics.Texture  // Texture handle
	magFilter    uint32            // magnification filter
	minFilter    uint32            // minification filter
	wrapS        uint32            // wrap mode for s coordinate
	wrapT        uint32            // wrap mode for t coordinate
	iformat      int32             // internal format
	width        int32             // texture width in pixels
	height       int32             // texture height in pixels
	format       uint32            // format of the pixel data
	formatType   uint32            // type of the pixel data
	updateData   bool              // texture data needs to be sent
	updateParams bool              // texture parameters needs to be sent
	genMipmap    bool              // generate mipmaps flag
	// data         interface{}    // array with texture data
	// uTexture     gls.Uniform1i  // Texture unit uniform
	// uFlipY       gls.Uniform1i  // Flip Y coordinate flag uniform
	// uVisible     gls.Uniform1i  // Texture visible uniform
	// uOffset      gls.Uniform2f  // Texture offset uniform
	// uRepeat      gls.Uniform2f  // Texture repeat uniform
}

func NewTexture2D() *texture2D {
	return &texture2D{}
}

func Texture2DFromImage(file string) (*texture2D, error) {
	return nil, nil
}

func (t *texture2D) SetImage(file string) error {
	return nil
}

func Texture2DFromRGBA() (*texture2D, error) {
	return nil, nil
}

func (t *texture2D) SetRGBA(i *image.RGBA) error {
	return nil
}

func Texture2DFromData() (*texture2D, error) {
	return nil, nil
}

func (t *texture2D) SetData() error {
	return nil
}

func (t *texture2D) Initialize() {
	t.p = nil
	t.refCount = 1
	t.handle = 0
	t.magFilter = graphics.LINEAR
	t.minFilter = graphics.LINEAR
	t.wrapS = graphics.CLAMP_TO_EDGE
	t.wrapT = graphics.CLAMP_TO_EDGE
	t.updateData = false
	t.updateParams = true
	t.genMipmap = true
	//t.uTexture.Init("MatTexture")
	//t.uFlipY.Init("MatTexFlipY")
	//t.uVisible.Init("MatTexVisible")
	//t.uOffset.Init("MatTexOffset")
	//t.uRepeat.Init("MatTexRepeat")
	//t.uRepeat.Set(1, 1)
	//t.uOffset.Set(0, 0)
	//t.uVisible.Set(1)
	//t.uFlipY.Set(1)
}

func (t *texture2D) Close() {
	//
}

func (t *texture2D) Increment() {
	t.refCount++
}

func (t *texture2D) Decrement() {
	t.refCount--
}

type Renderer interface {
	Render(graphics.Provider, int)
}

func (t *texture2D) Render(p graphics.Provider, idx int) {
	if t.p == nil {
		t.handle = p.GenTexture()
		t.p = p
	}

	if t.updateData {
		// Sets the texture unit for this texture
		p.ActiveTexture(graphics.Texture(graphics.TEXTURE0 + idx))
		p.BindTexture(graphics.TEXTURE_2D, t.handle)
		//p.TexImage2D(
		//	graphics.TEXTURE_2D, // texture type
		//	0,                   // level of detail
		//	t.iformat,           // internal format
		//	t.width,             // width in texels
		//	t.height,            // height in texels
		//	0,                   // border must be 0
		//	t.format,            // format of supplied texture data
		//	t.formatType,        // type of external format color component
		//	//t.data,              // image data
		//)
		// Generates mipmaps if requested
		if t.genMipmap {
			p.GenerateMipmap(graphics.TEXTURE_2D)
		}
		// No data to send
		t.updateData = false
	}

	// Sets the texture unit for this texture
	p.ActiveTexture(graphics.Texture(graphics.TEXTURE0 + idx))
	p.BindTexture(graphics.TEXTURE_2D, t.handle)

	// Sets texture parameters if needed
	//if t.updateParams {
	//	gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_MAG_FILTER, int32(t.magFilter))
	//	gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_MIN_FILTER, int32(t.minFilter))
	//	gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_WRAP_S, int32(t.wrapS))
	//	gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_WRAP_T, int32(t.wrapT))
	//	t.updateParams = false
	//}

	// Transfer uniforms
	//t.uTexture.Set(int32(idx))
	//t.uTexture.TransferIdx(gs, idx)
	//t.uFlipY.TransferIdx(gs, idx)
	//t.uVisible.TransferIdx(gs, idx)
	//t.uOffset.TransferIdx(gs, idx)
	//t.uRepeat.TransferIdx(gs, idx)
}
