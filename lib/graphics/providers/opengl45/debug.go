package opengl45

import (
	"strings"
	"unsafe"

	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/go-gl/gl/v4.5-core/gl"
)

type OGL45DEBUG struct {
	tag            string
	major, minor   int
	currentProgram graphics.Program
	debug          bool
}

type glNoReturn func()

func checkGLErr() {
	err := gl.GetError()
	if err != graphics.NO_ERROR {
		panic(err)
	}
}

func (g *OGL45DEBUG) run(fn glNoReturn) {
	fn()
	checkGLErr()
}

func (g *OGL45DEBUG) Init() error {
	err := gl.Init()
	if err != nil {
		return InitFailure(g.debug, err)
	}
	return nil
}

func (g *OGL45DEBUG) Version() (string, int, int) {
	return g.tag, g.major, g.minor
}

// ActiveTexture selects the active texture unit
func (g *OGL45DEBUG) ActiveTexture(t graphics.Texture) {
	g.run(func() { gl.ActiveTexture(uint32(t)) })
}

// AttachShader attaches a shader object to a program object
func (g *OGL45DEBUG) AttachShader(p graphics.Program, s graphics.Shader) {
	g.run(func() { gl.AttachShader(uint32(p), uint32(s)) })
}

// BindBuffer binds a buffer to the OpenGL target specified by enum
func (g *OGL45DEBUG) BindBuffer(target graphics.Enum, b graphics.Buffer) {
	g.run(func() { gl.BindBuffer(uint32(target), uint32(b)) })
}

// BindFragDataLocation binds a user-defined varying out variable
// to a fragment shader color number
func (g *OGL45DEBUG) BindFragDataLocation(p graphics.Program, color uint32, name string) {
	// name has to be zero terminated for gl.Str()
	glName := name + "\x00"
	g.run(func() { gl.BindFragDataLocation(uint32(p), color, gl.Str(glName)) })
}

// BindFramebuffer binds a framebuffer to a framebuffer target
func (g *OGL45DEBUG) BindFramebuffer(target graphics.Enum, fb graphics.Buffer) {
	g.run(func() { gl.BindFramebuffer(uint32(target), uint32(fb)) })
}

// BindRenderbuffer binds a renderbuffer to a renderbuffer target
func (g *OGL45DEBUG) BindRenderbuffer(target graphics.Enum, renderbuffer graphics.Buffer) {
	g.run(func() { gl.BindRenderbuffer(uint32(target), uint32(renderbuffer)) })
}

// BindTexture binds a texture to the OpenGL target specified by enum
func (g *OGL45DEBUG) BindTexture(target graphics.Enum, t graphics.Texture) {
	g.run(func() { gl.BindTexture(uint32(target), uint32(t)) })
}

// BindVertexArray binds a vertex array object
func (g *OGL45DEBUG) BindVertexArray(a uint32) {
	g.run(func() { gl.BindVertexArray(a) })
}

// BlendEquation specifies the equation used for both the RGB and
// alpha blend equations
func (g *OGL45DEBUG) BlendEquation(mode graphics.Enum) {
	g.run(func() { gl.BlendEquation(uint32(mode)) })
}

func (g *OGL45DEBUG) BlendEquationSeparate(modeRGB, modeAlpha graphics.Enum) {
	g.run(func() { gl.BlendEquationSeparate(uint32(modeRGB), uint32(modeAlpha)) })
}

// BlendFunc specifies the pixel arithmetic for the blend fucntion
func (g *OGL45DEBUG) BlendFunc(sFactor, dFactor graphics.Enum) {
	g.run(func() { gl.BlendFunc(uint32(sFactor), uint32(dFactor)) })
}

func (g *OGL45DEBUG) BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha graphics.Enum) {
	g.run(
		func() {
			gl.BlendFuncSeparate(
				uint32(sfactorRGB),
				uint32(dfactorRGB),
				uint32(sfactorAlpha),
				uint32(dfactorAlpha),
			)
		})
}

// BlitFramebuffer copies a block of pixels from one framebuffer object to another
func (g *OGL45DEBUG) BlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1 int32, mask graphics.Bitfield, filter graphics.Enum) {
	g.run(func() {
		gl.BlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1, uint32(mask), uint32(filter))
	})
}

// BufferData creates a new data store for the bound buffer object.
func (g *OGL45DEBUG) BufferData(target graphics.Enum, size int, data unsafe.Pointer, usage graphics.Enum) {
	g.run(func() { gl.BufferData(uint32(target), size, data, uint32(usage)) })
}

// CheckFramebufferStatus checks the completeness status of a framebuffer
func (g *OGL45DEBUG) CheckFramebufferStatus(target graphics.Enum) graphics.Enum {
	ret := graphics.Enum(gl.CheckFramebufferStatus(uint32(target)))
	checkGLErr()
	return ret
}

// Clear clears the window buffer specified in mask
func (g *OGL45DEBUG) Clear(mask graphics.Enum) {
	g.run(func() { gl.Clear(uint32(mask)) })
}

// ClearColor specifies the RGBA value used to clear the color buffers
func (g *OGL45DEBUG) ClearColor(red, green, blue, alpha float32) {
	g.run(func() { gl.ClearColor(red, green, blue, alpha) })
}

//
func (g *OGL45DEBUG) ClearDepth(v float32) {
	g.run(func() { gl.ClearDepthf(v) })
}

//
func (g *OGL45DEBUG) ClearStencil(v int32) {
	g.run(func() { gl.ClearStencil(v) })
}

// CompileShader compiles the shader object
func (g *OGL45DEBUG) CompileShader(s graphics.Shader) {
	g.run(func() { gl.CompileShader(uint32(s)) })
}

// CreateProgram creates a new shader program object
func (g *OGL45DEBUG) CreateProgram() graphics.Program {
	prog := graphics.Program(gl.CreateProgram())
	checkGLErr()
	return prog
}

// CreateShader creates a new shader object
func (g *OGL45DEBUG) CreateShader(ty graphics.Enum) graphics.Shader {
	shdr := graphics.Shader(gl.CreateShader(uint32(ty)))
	checkGLErr()
	return shdr
}

// CullFace specifies whether to use front or back face culling
func (g *OGL45DEBUG) CullFace(mode graphics.Enum) {
	g.run(func() { gl.CullFace(uint32(mode)) })
}

func (g *OGL45DEBUG) CurrentProgram() graphics.Program {
	return g.currentProgram
}

// DeleteBuffer deletes the OpenGL buffer object
func (g *OGL45DEBUG) DeleteBuffer(b graphics.Buffer) {
	uintV := uint32(b)
	g.run(func() { gl.DeleteBuffers(1, &uintV) })
}

// DeleteFramebuffer deletes the framebuffer object
func (g *OGL45DEBUG) DeleteFramebuffer(fb graphics.Buffer) {
	uintV := uint32(fb)
	g.run(func() { gl.DeleteFramebuffers(1, &uintV) })
}

// DeleteProgram deletes the shader program object
func (g *OGL45DEBUG) DeleteProgram(p graphics.Program) {
	g.run(func() { gl.DeleteProgram(uint32(p)) })
	//g.currentProgram = 0
}

// DeleteRenderbuffer deletes the renderbuffer object
func (g *OGL45DEBUG) DeleteRenderbuffer(rb graphics.Buffer) {
	uintV := uint32(rb)
	g.run(func() { gl.DeleteRenderbuffers(1, &uintV) })
}

// DeleteShader deletes the shader object
func (g *OGL45DEBUG) DeleteShader(s graphics.Shader) {
	g.run(func() { gl.DeleteShader(uint32(s)) })
}

// DeleteTexture deletes the specified texture
func (g *OGL45DEBUG) DeleteTexture(v graphics.Texture) {
	uintV := uint32(v)
	g.run(func() { gl.DeleteTextures(1, &uintV) })
}

// DeleteVertexArray deletes an OpenGL VAO
func (g *OGL45DEBUG) DeleteVertexArray(a uint32) {
	uintV := uint32(a)
	g.run(func() { gl.DeleteVertexArrays(1, &uintV) })
}

//
func (g *OGL45DEBUG) DepthFunc(e graphics.Enum) {
	g.run(func() { gl.DepthFunc(uint32(e)) })
}

// DepthMask enables or disables writing into the depth buffer
func (g *OGL45DEBUG) DepthMask(flag bool) {
	g.run(func() { gl.DepthMask(flag) })
}

// Disable disables various GL capabilities.
func (g *OGL45DEBUG) Disable(e graphics.Enum) {
	g.run(func() { gl.Disable(uint32(e)) })
}

// DrawBuffers specifies a list of color buffers to be drawn into
func (g *OGL45DEBUG) DrawBuffers(buffers []uint32) {
	c := int32(len(buffers))
	g.run(func() { gl.DrawBuffers(c, &buffers[0]) })
}

// DrawElements renders primitives from array data
func (g *OGL45DEBUG) DrawElements(mode graphics.Enum, count int32, ty graphics.Enum, indices unsafe.Pointer) {
	g.run(func() { gl.DrawElements(uint32(mode), count, uint32(ty), indices) })
}

// DrawArrays renders primitives from array data
func (g *OGL45DEBUG) DrawArrays(mode graphics.Enum, first int32, count int32) {
	g.run(func() { gl.DrawArrays(uint32(mode), first, count) })
}

// Enable enables various GL capabilities.
func (g *OGL45DEBUG) Enable(e graphics.Enum) {
	g.run(func() { gl.Enable(uint32(e)) })
}

// EnableVertexAttribArray enables a vertex attribute array
func (g *OGL45DEBUG) EnableVertexAttribArray(a uint32) {
	g.run(func() { gl.EnableVertexAttribArray(a) })
}

// FramebufferRenderbuffer attaches a renderbuffer as a logical buffer
// of a framebuffer object
func (g *OGL45DEBUG) FramebufferRenderbuffer(target, attachment, renderbuffertarget graphics.Enum, renderbuffer graphics.Buffer) {
	g.run(func() {
		gl.FramebufferRenderbuffer(uint32(target), uint32(attachment), uint32(renderbuffertarget), uint32(renderbuffer))
	})
}

// FramebufferTexture2D attaches a texture object to a framebuffer
func (g *OGL45DEBUG) FramebufferTexture2D(target, attachment, textarget graphics.Enum, texture graphics.Texture, level int32) {
	g.run(func() {
		gl.FramebufferTexture2D(uint32(target), uint32(attachment), uint32(textarget), uint32(texture), level)
	})
}

//
func (g *OGL45DEBUG) FrontFace(e graphics.Enum) {
	g.run(func() { gl.FrontFace(uint32(e)) })
}

// GenBuffer creates an OpenGL buffer object
func (g *OGL45DEBUG) GenBuffer() graphics.Buffer {
	var b uint32
	gl.GenBuffers(1, &b)
	checkGLErr()
	return graphics.Buffer(b)
}

// GenerateMipmap generates mipmaps for a specified texture target
func (g *OGL45DEBUG) GenerateMipmap(t graphics.Enum) {
	g.run(func() { gl.GenerateMipmap(uint32(t)) })
}

// GenFramebuffer generates a OpenGL framebuffer object
func (g *OGL45DEBUG) GenFramebuffer() graphics.Buffer {
	var b uint32
	gl.GenFramebuffers(1, &b)
	checkGLErr()
	return graphics.Buffer(b)
}

// GenRenderbuffer generates a OpenGL renderbuffer object
func (g *OGL45DEBUG) GenRenderbuffer() graphics.Buffer {
	var b uint32
	gl.GenRenderbuffers(1, &b)
	checkGLErr()
	return graphics.Buffer(b)
}

// GenTexture creates an OpenGL texture object
func (g *OGL45DEBUG) GenTexture() graphics.Texture {
	var t uint32
	gl.GenTextures(1, &t)
	checkGLErr()
	return graphics.Texture(t)
}

// GenVertexArray creates an OpenGL VAO
func (g *OGL45DEBUG) GenVertexArray() uint32 {
	var a uint32
	gl.GenVertexArrays(1, &a)
	checkGLErr()
	return a
}

// GetAttribLocation returns the location of a attribute variable
func (g *OGL45DEBUG) GetAttribLocation(p graphics.Program, name string) int32 {
	glName := name + "\x00"
	ret := gl.GetAttribLocation(uint32(p), gl.Str(glName))
	checkGLErr()
	return ret
}

//
func (g *OGL45DEBUG) GetAttribCurrentLocation(name string) int32 {
	glName := name + "\x00"
	ret := gl.GetAttribLocation(uint32(g.currentProgram), gl.Str(glName))
	checkGLErr()
	return ret
}

// GetError returns the next error
func (g *OGL45DEBUG) GetError() uint32 {
	return gl.GetError()
}

// GetProgramInfoLog returns the information log for a program object
func (g *OGL45DEBUG) GetProgramInfoLog(p graphics.Program) string {
	var logLength int32
	g.GetProgramiv(p, gl.INFO_LOG_LENGTH, &logLength)

	// make sure the string is zero'd out to start with
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(uint32(p), logLength, nil, gl.Str(log))

	return log
}

// GetProgramiv returns a parameter from the program object
func (g *OGL45DEBUG) GetProgramiv(p graphics.Program, pname graphics.Enum, params *int32) {
	g.run(func() { gl.GetProgramiv(uint32(p), uint32(pname), params) })
}

// GetShaderInfoLog returns the information log for a shader object
func (g *OGL45DEBUG) GetShaderInfoLog(s graphics.Shader) string {
	var logLength int32
	g.GetShaderiv(s, gl.INFO_LOG_LENGTH, &logLength)

	// make sure the string is zero'd out to start with
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(uint32(s), logLength, nil, gl.Str(log))

	return log
}

// GetShaderiv returns a parameter from the shader object
func (g *OGL45DEBUG) GetShaderiv(s graphics.Shader, pname graphics.Enum, params *int32) {
	g.run(func() { gl.GetShaderiv(uint32(s), uint32(pname), params) })
}

// GetUniformLocation returns the location of a uniform variable
func (g *OGL45DEBUG) GetUniformLocation(p graphics.Program, name string) int32 {
	glName := name + "\x00"
	ret := gl.GetUniformLocation(uint32(p), gl.Str(glName))
	checkGLErr()
	return ret
}

func (g *OGL45DEBUG) GetUniformCurrentLocation(name string) int32 {
	glName := name + "\x00"
	ret := gl.GetUniformLocation(uint32(g.currentProgram), gl.Str(glName))
	checkGLErr()
	return ret
}

func (g *OGL45DEBUG) LineWidth(w float32) {
	g.run(func() { gl.LineWidth(w) })
}

// LinkProgram links a program object
func (g *OGL45DEBUG) LinkProgram(p graphics.Program) {
	g.run(func() { gl.LinkProgram(uint32(p)) })
}

// PolygonMode sets a polygon rasterization mode.
func (g *OGL45DEBUG) PolygonMode(face, mode graphics.Enum) {
	g.run(func() { gl.PolygonMode(uint32(face), uint32(mode)) })
}

// PolygonOffset sets the scale and units used to calculate depth values
func (g *OGL45DEBUG) PolygonOffset(factor float32, units float32) {
	g.run(func() { gl.PolygonOffset(factor, units) })
}

// Ptr takes a slice or a pointer and returns an OpenGL compatbile address
func (g *OGL45DEBUG) Ptr(data interface{}) unsafe.Pointer {
	return gl.Ptr(data)
}

// PtrOffset takes a pointer offset and returns a GL-compatible pointer.
// Useful for functions such as glVertexAttribPointer that take pointer
// parameters indicating an offset rather than an absolute memory address.
func (g *OGL45DEBUG) PtrOffset(offset int) unsafe.Pointer {
	return gl.PtrOffset(offset)
}

// ReadBuffer specifies the color buffer source for pixels
func (g *OGL45DEBUG) ReadBuffer(src graphics.Enum) {
	g.run(func() { gl.ReadBuffer(uint32(src)) })
}

// RenderbufferStorage establishes the format and dimensions of a renderbuffer
func (g *OGL45DEBUG) RenderbufferStorage(target graphics.Enum, internalformat graphics.Enum, width int32, height int32) {
	g.run(func() { gl.RenderbufferStorage(uint32(target), uint32(internalformat), width, height) })
}

// RenderbufferStorageMultisample establishes the format and dimensions of a renderbuffer
func (g *OGL45DEBUG) RenderbufferStorageMultisample(target graphics.Enum, samples int32, internalformat graphics.Enum, width int32, height int32) {
	g.run(
		func() {
			gl.RenderbufferStorageMultisample(
				uint32(target),
				samples,
				uint32(internalformat),
				width,
				height,
			)
		},
	)
}

// Scissor clips to a rectangle with the location and dimensions specified.
func (g *OGL45DEBUG) Scissor(x, y, w, h int32) {
	g.run(func() { gl.Scissor(x, y, w, h) })
}

// ShaderSource replaces the source code for a shader object.
func (g *OGL45DEBUG) ShaderSource(s graphics.Shader, source string) {
	glSource, free := gl.Strs(source + "\x00")
	g.run(func() { gl.ShaderSource(uint32(s), 1, glSource, nil) })
	free()
}

// TexImage2D writes a 2D texture image.
func (g *OGL45DEBUG) TexImage2D(target graphics.Enum, level, intfmt, width, height, border int32, format graphics.Enum, ty graphics.Enum, ptr unsafe.Pointer, dataLength int) {
	g.run(func() {
		gl.TexImage2D(uint32(target), level, intfmt, width, height, border, uint32(format), uint32(ty), ptr)
	})
}

// TexImage2DMultisample establishes the data storage, format, dimensions, and number of samples of a multisample texture's image
func (g *OGL45DEBUG) TexImage2DMultisample(target graphics.Enum, samples int32, intfmt graphics.Enum, width int32, height int32, fixedsamplelocations bool) {
	g.run(func() {
		gl.TexImage2DMultisample(uint32(target), samples, uint32(intfmt), width, height, fixedsamplelocations)
	})
}

// TexParameterf sets a float texture parameter
func (g *OGL45DEBUG) TexParameterf(target, pname graphics.Enum, param float32) {
	g.run(func() { gl.TexParameterf(uint32(target), uint32(pname), param) })
}

// TexParameterfv sets a float texture parameter
func (g *OGL45DEBUG) TexParameterfv(target, pname graphics.Enum, params *float32) {
	g.run(func() { gl.TexParameterfv(uint32(target), uint32(pname), params) })
}

// TexParameteri sets a float texture parameter
func (g *OGL45DEBUG) TexParameteri(target, pname graphics.Enum, param int32) {
	g.run(func() { gl.TexParameteri(uint32(target), uint32(pname), param) })
}

// TexStorage3D simultaneously specifies storage for all levels of a three-dimensional,
// two-dimensional array or cube-map array texture
func (g *OGL45DEBUG) TexStorage3D(target graphics.Enum, level int32, intfmt uint32, width, height, depth int32) {
	g.run(func() { gl.TexStorage3D(uint32(target), level, intfmt, width, height, depth) })
}

// TexSubImage3D specifies a three-dimensonal texture subimage
func (g *OGL45DEBUG) TexSubImage3D(target graphics.Enum, level, xoff, yoff, zoff, width, height, depth int32, fmt, ty graphics.Enum, ptr unsafe.Pointer) {
	g.run(func() {
		gl.TexSubImage3D(uint32(target), level, xoff, yoff, zoff, width, height, depth, uint32(fmt), uint32(ty), ptr)
	})
}

// Uniform1i specifies the value of a uniform variable for the current program object
func (g *OGL45DEBUG) Uniform1i(location int32, v int32) {
	g.run(func() { gl.Uniform1i(location, v) })
}

// Uniform1iv specifies the value of a uniform variable for the current program object
func (g *OGL45DEBUG) Uniform1iv(location int32, values []int32) {
	g.run(func() { gl.Uniform1iv(location, int32(len(values)), &values[0]) })
}

// Uniform1f specifies the value of a uniform variable for the current program object
func (g *OGL45DEBUG) Uniform1f(location int32, v float32) {
	g.run(func() { gl.Uniform1f(location, v) })
}

// Uniform1fv specifies the value of a uniform variable for the current program object
func (g *OGL45DEBUG) Uniform1fv(location int32, values []float32) {
	g.run(func() { gl.Uniform1fv(location, int32(len(values)), &values[0]) })
}

// Uniform3f specifies the value of a uniform variable for the current program object
func (g *OGL45DEBUG) Uniform3f(location int32, v0, v1, v2 float32) {
	g.run(func() { gl.Uniform3f(location, v0, v1, v2) })
}

// Uniform3fv specifies the value of a uniform variable for the current program object
func (g *OGL45DEBUG) Uniform3fv(location int32, values []float32) {
	g.run(func() { gl.Uniform3fv(location, int32(len(values)), &values[0]) })
}

// Uniform4f specifies the value of a uniform variable for the current program object
func (g *OGL45DEBUG) Uniform4f(location int32, v0, v1, v2, v3 float32) {
	g.run(func() { gl.Uniform4f(location, v0, v1, v2, v3) })
}

// Uniform4fv specifies the value of a uniform variable for the current program object
func (g *OGL45DEBUG) Uniform4fv(location int32, values []float32) {
	g.run(func() { gl.Uniform4fv(location, int32(len(values)), &values[0]) })
}

// UniformMatrix4fv specifies the value of a uniform variable for the current program object
func (g *OGL45DEBUG) UniformMatrix4fv(location, count int32, transpose bool, value []float32) {
	g.run(func() { gl.UniformMatrix4fv(location, count, transpose, &value[0]) })
}

// UseProgram installs a program object as part of the current rendering state
func (g *OGL45DEBUG) UseProgram(p graphics.Program) {
	g.currentProgram = p
	g.run(func() { gl.UseProgram(uint32(p)) })
}

// VertexAttribPointer uses a bound buffer to define vertex attribute data.
//
// The size argument specifies the number of components per attribute,
// between 1-4. The stride argument specifies the byte offset between
// consecutive vertex attributes.
func (g *OGL45DEBUG) VertexAttribPointer(dst uint32, size int32, ty graphics.Enum, normalized bool, stride int32, ptr unsafe.Pointer) {
	g.run(func() { gl.VertexAttribPointer(dst, size, uint32(ty), normalized, stride, ptr) })
}

// VertexAttribPointer uses a bound buffer to define vertex attribute data.
// Only integer types are accepted by this function.
func (g *OGL45DEBUG) VertexAttribIPointer(dst uint32, size int32, ty graphics.Enum, stride int32, ptr unsafe.Pointer) {
	g.run(func() { gl.VertexAttribIPointer(dst, size, uint32(ty), stride, ptr) })
}

// Viewport sets the viewport, an affine transformation that
// normalizes device coordinates to window coordinates.
func (g *OGL45DEBUG) Viewport(x, y, width, height int32) {
	g.run(func() { gl.Viewport(x, y, width, height) })
}
