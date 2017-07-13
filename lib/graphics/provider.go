package graphics

import (
	"unsafe"

	"github.com/Laughs-In-Flowers/shiva/lib/xrror"
)

type ProviderT int

func (t ProviderT) String() string {
	switch t {
	case OPENGL45:
		return "opengl4.5"
	}
	return "UNKNOWN"
}

func StringToProviderT(s string) ProviderT {
	switch s {
	case "opengl4.5":
		return OPENGL45
	}
	return UNKNOWN
}

const (
	OPENGL45 ProviderT = iota
	UNKNOWN
)

var DefaultProvider = OPENGL45

var UnknownGraphicsImplementation = xrror.Xrror("%s is not a known graphics implementation").Out

type NewProviderFunc func(bool) Provider

type gpr struct {
	has map[ProviderT]NewProviderFunc
}

func (g *gpr) add(name string, fn NewProviderFunc) {
	pn := StringToProviderT(name)
	g.has[pn] = fn
}

func (g *gpr) get(name string, debug bool) (Provider, error) {
	var p Provider
	pt := StringToProviderT(name)
	if pfn, exists := g.has[pt]; exists {
		p = pfn(debug)
		return p, nil
	}
	return nil, UnknownGraphicsImplementation(name)
}

var ProviderRegistry *gpr

func Register(name string, fn NewProviderFunc) {
	ProviderRegistry.add(name, fn)
}

func New(s string, debug bool) (Provider, error) {
	return ProviderRegistry.get(s, debug)
}

// Enum is a type indicating the uint32 use as an enumeration value
type Enum uint32

// Texture is a type indicating the uint32 use as a graphics level texture
type Texture uint32

// Buffer is a type indicating the uint32 use as a graphics level buffer
type Buffer uint32

// Program is a type indicating the uint32 use as a graphics level shader program
type Program uint32

// Shader is a type indicating the uint32 use as a graphics level shader
type Shader uint32

// Bitfield is a typ indicating the uint32 use as a graphics level bitfield
type Bitfield uint32

// GraphicsProvider represents a common way to interface with graphics
// 'drivers' like OpenGL or OpenGL ES.
type Provider interface {
	// Return version information as string followed by two integers.
	Version() (string, int, int)

	// Initialize the provider, returning an error message.
	Init() error

	// ActiveTexture selects the active texture unit
	ActiveTexture(Texture)

	// AttachShader attaches a shader object to a program object
	AttachShader(Program, Shader)

	// BindBuffer binds a buffer to the graphics level target specified by enum
	BindBuffer(Enum, Buffer)

	// BindFragDataLocation binds a user-defined varying out variable
	// to a fragment shader color number
	BindFragDataLocation(Program, uint32, string)

	// BindFramebuffer binds a framebuffer to a framebuffer target
	BindFramebuffer(Enum, Buffer)

	// BindRenderbuffer binds a renderbuffer to a renderbuffer target
	BindRenderbuffer(Enum, Buffer)

	// BindTexture binds a texture to the graphics level target specified by enum
	BindTexture(Enum, Texture)

	// BindVertexArray binds a vertex array object
	BindVertexArray(uint32)

	// BlendEquation specifies the equation used for both the RGB and
	// alpha blend equations
	BlendEquation(Enum)

	// set the RGB blend equation and the alpha blend equation separately
	BlendEquationSeparate(Enum, Enum)

	// BlendFunc specifies the pixel arithmetic for the blend fucntion
	BlendFunc(Enum, Enum)

	//
	BlendFuncSeparate(Enum, Enum, Enum, Enum)

	// BlitFramebuffer copies a block of pixels from one framebuffer object to another
	BlitFramebuffer(int32, int32, int32, int32, int32, int32, int32, int32, Bitfield, Enum)

	// BufferData creates a new data store for the bound buffer object.
	BufferData(Enum, int, unsafe.Pointer, Enum)

	// CheckFramebufferStatus checks the completeness status of a framebuffer
	CheckFramebufferStatus(Enum) Enum

	// Clear clears the window buffer specified in mask
	Clear(Enum)

	// ClearColor specifies the RGBA value used to clear the color buffers
	ClearColor(float32, float32, float32, float32)

	// CompileShader compiles the shader object
	CompileShader(Shader)

	// CreateProgram creates a new shader program object
	CreateProgram() Program

	// CreateShader creates a new shader object
	CreateShader(Enum) Shader

	// CullFace specifies whether to use front or back face culling
	CullFace(Enum)

	// CurrentProgram will return the currently bound shader program, 0 if no program is bound.
	CurrentProgram() Program

	// DeleteBuffer deletes the graphics level buffer object
	DeleteBuffer(Buffer)

	// DeleteFramebuffer deletes the framebuffer object
	DeleteFramebuffer(Buffer)

	// DeleteProgram deletes the shader program object
	DeleteProgram(Program)

	// DeleteRenderbuffer deletes the renderbuffer object
	DeleteRenderbuffer(Buffer)

	// DeleteShader deletes the shader object
	DeleteShader(Shader)

	// DeleteTexture deletes the specified texture
	DeleteTexture(Texture)

	// DeleteVertexArray deletes a graphics level VAO
	DeleteVertexArray(uint32)

	//
	DepthFunc(Enum)

	// DepthMask enables or disables writing into the depth buffer
	DepthMask(bool)

	// Disable disables various graphics level capabilities
	Disable(Enum)

	// DrawBuffers specifies a list of color buffers to be drawn into
	DrawBuffers([]uint32)

	// DrawElements renders primitives from array data
	DrawElements(Enum, int32, Enum, unsafe.Pointer)

	// DrawArrays renders primitives from array data
	DrawArrays(Enum, int32, int32)

	// Enable enables various graphics level capabilities.
	Enable(Enum)

	// EnableVertexAttribArray enables a vertex attribute array
	EnableVertexAttribArray(uint32)

	// FramebufferRenderbuffer attaches a renderbuffer as a logical buffer
	// of a framebuffer object
	FramebufferRenderbuffer(Enum, Enum, Enum, Buffer)

	// FramebufferTexture2D attaches a texture object to a framebuffer
	FramebufferTexture2D(Enum, Enum, Enum, Texture, int32)

	//
	FrontFace(Enum)

	// GenBuffer creates a graphics level buffer object
	GenBuffer() Buffer

	// GenerateMipmap generates mipmaps for a specified texture target
	GenerateMipmap(Enum)

	// GenFramebuffer generates a graphics level framebuffer object
	GenFramebuffer() Buffer

	// GenRenderbuffer generates a graphics level renderbuffer object
	GenRenderbuffer() Buffer

	// GenTexture creates a graphics level texture object
	GenTexture() Texture

	// GenVertexArray creates a graphics level VAO
	GenVertexArray() uint32

	// GetAttribLocation returns the location of a attribute variable
	GetAttribLocation(Program, string) int32

	// GetCurrentUniformLocation returns the location of a uniform variable relative to current program
	GetCurrentUniformLocation(string) int32

	// GetError returns the next error
	GetError() uint32

	// GetProgramInfoLog returns the information log for a program object
	GetProgramInfoLog(Program) string

	// GetProgramiv returns a parameter from the program object
	GetProgramiv(Program, Enum, *int32)

	// GetShaderInfoLog returns the information log for a shader object
	GetShaderInfoLog(Shader) string

	// GetShaderiv returns a parameter from the shader object
	GetShaderiv(Shader, Enum, *int32)

	// GetUniformLocation returns the location of a uniform variable
	GetUniformLocation(Program, string) int32

	//
	LineWidth(float32)

	// LinkProgram links a program object
	LinkProgram(Program)

	// PolygonMode sets a polygon rasterization mode.
	PolygonMode(Enum, Enum)

	// PolygonOffset sets the scale and units used to calculate depth values
	PolygonOffset(float32, float32)

	// Ptr takes a slice or a pointer and returns a graphics level compatbile address
	Ptr(interface{}) unsafe.Pointer

	// PtrOffset takes a pointer offset and returns a GL-compatible pointer.
	// Useful for functions such as glVertexAttribPointer that take pointer
	// parameters indicating an offset rather than an absolute memory address.
	PtrOffset(int) unsafe.Pointer

	// ReadBuffer specifies the color buffer source for pixels
	ReadBuffer(Enum)

	// RenderbufferStorage establishes the format and dimensions of a renderbuffer
	RenderbufferStorage(Enum, Enum, int32, int32)

	// RenderbufferStorageMultisample establishes the format and dimensions of a renderbuffer
	RenderbufferStorageMultisample(Enum, int32, Enum, int32, int32)

	// Scissor clips to a rectangle with the location and dimensions specified.
	Scissor(int32, int32, int32, int32)

	// ShaderSource replaces the source code for a shader object.
	ShaderSource(Shader, string)

	// TexImage2D writes a 2D texture image.
	TexImage2D(Enum, int32, int32, int32, int32, int32, Enum, Enum, unsafe.Pointer, int)

	// TexImage2DMultisample establishes the data storage, format, dimensions, and number of samples of a multisample texture's image
	TexImage2DMultisample(Enum, int32, Enum, int32, int32, bool)

	// TexParameterf sets a float texture parameter
	TexParameterf(Enum, Enum, float32)

	// TexParameterfv sets a float texture parameter
	TexParameterfv(Enum, Enum, *float32)

	// TexParameteri sets an int texture parameter
	TexParameteri(Enum, Enum, int32)

	// TexStorage3D simultaneously specifies storage for all levels of a three-dimensional,
	// two-dimensional array or cube-map array texture
	TexStorage3D(Enum, int32, uint32, int32, int32, int32)

	// TexSubImage3D specifies a three-dimensonal texture subimage
	TexSubImage3D(Enum, int32, int32, int32, int32, int32, int32, int32, Enum, Enum, unsafe.Pointer)

	// Uniform1i specifies the value of a uniform variable for the current program object
	Uniform1i(int32, int32)

	// Uniform1iv specifies the value of a uniform variable for the current program object
	Uniform1iv(int32, []int32)

	// Uniform1f specifies the value of a uniform variable for the current program object
	Uniform1f(int32, float32)

	// Uniform1fv specifies the value of a uniform variable for the current program object
	Uniform1fv(int32, []float32)

	// Uniform3f specifies the value of a uniform variable for the current program object
	Uniform3f(int32, float32, float32, float32)

	// Uniform3fv specifies the value of a uniform variable for the current program object
	Uniform3fv(int32, []float32)

	// UniformMatrix3fv specifies the value of a uniform variable for the current program object
	//UniformMatrix3fv(int32, int32, bool, []float32)

	// Uniform4f specifies the value of a uniform variable for the current program object
	Uniform4f(int32, float32, float32, float32, float32)

	// Uniform4fv specifies the value of a uniform variable for the current program object
	Uniform4fv(int32, []float32)

	// UniformMatrix4fv specifies the value of a uniform variable for the current program object
	UniformMatrix4fv(int32, int32, bool, []float32)

	// UseProgram installs a program object as part of the current rendering state
	UseProgram(Program)

	// VertexAttribPointer uses a bound buffer to define vertex attribute data.
	//
	// The size argument specifies the number of components per attribute,
	// between 1-4. The stride argument specifies the byte offset between
	// consecutive vertex attributes.
	VertexAttribPointer(uint32, int32, Enum, bool, int32, unsafe.Pointer)

	// VertexAttribIPointer uses a bound buffer to define vertex attribute data.
	// Only integer types are accepted by this function.
	VertexAttribIPointer(uint32, int32, Enum, int32, unsafe.Pointer)

	// Viewport sets the viewport, an affine transformation that
	// normalizes device coordinates to window coordinates.
	Viewport(int32, int32, int32, int32)
}

func init() {
	ProviderRegistry = &gpr{make(map[ProviderT]NewProviderFunc)}
}
