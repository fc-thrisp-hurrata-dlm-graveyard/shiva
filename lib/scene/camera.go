package scene

import (
	"strings"

	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/Laughs-In-Flowers/shiva/lib/math"
	"github.com/Laughs-In-Flowers/shiva/lib/render"
	"github.com/Laughs-In-Flowers/shiva/lib/xrror"

	l "github.com/yuin/gopher-lua"
)

type CamT int

const (
	NOT_A_CAMERA_TYPE CamT = iota
	PERSPECTIVE
	ORTHOGRAPHIC
	CUSTOM
)

func StringToCamT(s string) CamT {
	switch strings.ToLower(s) {
	case "perspective":
		return PERSPECTIVE
	case "orthographic":
		return ORTHOGRAPHIC
	case "custom":
		return CUSTOM
	}
	return NOT_A_CAMERA_TYPE
}

type Plane int

const (
	FOV Plane = iota
	ASPECT
	NEAR
	FAR
	ZOOM
	LEFT
	RIGHT
	TOP
	BOTTOM
	UNKNOWN_PLANE
)

func StringToPlane(s string) Plane {
	switch strings.ToLower(s) {
	case "fov":
		return FOV
	case "aspect":
		return ASPECT
	case "near":
		return NEAR
	case "far":
		return FAR
	case "zoom":
		return ZOOM
	case "left":
		return LEFT
	case "right":
		return RIGHT
	case "top":
		return TOP
	case "bottom":
		return BOTTOM
	}
	return UNKNOWN_PLANE
}

type Cam interface {
	Planes() []float32
	GetPlane(Plane) float32
	SetPlane(Plane, float32)
	ProjectionMatrix() math.Matrice
}

type planesFunc func(*cam) []float32

type projectionMatrixFunc func(*cam)

type cam struct {
	camt                         CamT
	plfn                         planesFunc
	fov, aspect, near, far, zoom float32
	left, right, top, bottom     float32
	pmfn                         projectionMatrixFunc
	projectionChange             bool
	projectionMatrix             math.Matrice
}

func newCam(t CamT, p planesFunc, m projectionMatrixFunc) *cam {
	return &cam{
		camt:             t,
		plfn:             p,
		zoom:             1.0,
		pmfn:             m,
		projectionChange: true,
		projectionMatrix: math.Mat4(),
	}
}

func (c *cam) Planes() []float32 {
	return c.plfn(c)
}

func (c *cam) GetPlane(p Plane) float32 {
	var ret float32
	switch p {
	case FOV:
		ret = c.fov
	case ASPECT:
		ret = c.aspect
	case NEAR:
		ret = c.near
	case FAR:
		ret = c.far
	case ZOOM:
		ret = c.zoom
	case LEFT:
		ret = c.left
	case RIGHT:
		ret = c.right
	case TOP:
		ret = c.top
	case BOTTOM:
		ret = c.bottom
	}
	return ret
}

func (c *cam) SetPlane(p Plane, v float32) {
	switch p {
	case FOV:
		c.projectionChange = true
		c.fov = v
	case ASPECT:
		c.projectionChange = true
		c.aspect = v
	case NEAR:
		c.near = v
	case FAR:
		c.far = v
	case ZOOM:
		c.projectionChange = true
		c.zoom = math.Abs(v)
	case LEFT:
		c.left = v
	case RIGHT:
		c.right = v
	case TOP:
		c.top = v
	case BOTTOM:
		c.bottom = v
	}
}

func (c *cam) ProjectionMatrix() math.Matrice {
	if c.projectionChange {
		c.pmfn(c)
		c.projectionChange = false
	}
	return c.projectionMatrix
}

func perspective(fov, aspect, near, far float32) *cam {
	c := newCam(
		PERSPECTIVE,
		func(c *cam) []float32 {
			return []float32{
				c.GetPlane(FOV),
				c.GetPlane(ASPECT),
				c.GetPlane(NEAR),
				c.GetPlane(FAR),
			}
		},
		func(c *cam) {
			p := c.Planes()
			c.projectionMatrix.Perspective(p[0], p[1], p[2], p[3])
		},
	)
	c.fov = fov
	c.aspect = aspect
	c.near = near
	c.far = far
	return c
}

func orthographic(left, right, top, bottom, near, far float32) *cam {
	c := newCam(
		ORTHOGRAPHIC,
		func(c *cam) []float32 {
			return []float32{
				c.GetPlane(LEFT),
				c.GetPlane(RIGHT),
				c.GetPlane(TOP),
				c.GetPlane(BOTTOM),
				c.GetPlane(NEAR),
				c.GetPlane(FAR),
				c.GetPlane(ZOOM),
			}
		},
		func(c *cam) {
			p := c.Planes()
			z := p[6]
			c.projectionMatrix.Orthographic(
				p[0]/z,
				p[1]/z,
				p[2]/z,
				p[3]/z,
				p[4],
				p[5],
			)
		},
	)
	c.left = left
	c.right = right
	c.top = top
	c.bottom = bottom
	c.near = near
	c.far = far
	return c
}

type camera struct {
	*cam
	*position
	target     math.Vector
	up         math.Vector
	viewMatrix math.Matrice
}

func newCamera(ck *cam) *camera {
	c := &camera{
		ck,
		newPosition(),
		math.Vec3(0, 0, 0),
		math.Vec3(0, 1, 0),
		math.Mat4(),
	}
	c.Update(DIRECTION, 0, 0, -1)
	return c
}

// update camera quat on changes
func (c *camera) update() {
	dr := c.direct(math.Vec3(0, 0, 0))
	q := math.Quat(0, 0, 0, 0)
	math.SetQuatFromUnitVectors(q, math.Vec3(0, 0, -1), dr)
	c.Update(ROTATE, q.Raw()...)
}

func (c *camera) Direction() math.Vector {
	tr := c.translate(math.Vec3(0, 0, 0))
	res := c.target
	res.Sub(tr).Normalize()
	return res
}

func (c *camera) LookAt(t math.Vector) {
	c.target = t
}

func (c *camera) ViewMatrix() math.Matrice {
	tr := c.translate(math.Vec3(0, 0, 0))
	c.viewMatrix.LookAt(tr, c.target, c.up)
	return c.viewMatrix
}

type cameraNode struct {
	*camera
	*node
}

const lCameraNodeClass = "NCAMERA"

func Camera(tag string, c *cam) Node {
	cc := newCamera(c)
	nn := newNode(tag, func(r render.Renderer, n Node) {
		cc.updateMatrixWorld(r)
		r.SetViewMatrice(cc.ViewMatrix())
		r.SetProjectionMatrice(cc.ProjectionMatrix())
	}, defaultRemovalFn, defaultReplaceFn, lCameraNodeClass, lNodeClass)

	return &cameraNode{
		cc,
		nn,
	}
}

var (
	cameraTag             TagFunc = tagFnFor("camera", 1)
	CameraFromStringError         = xrror.Xrror("%s is not a camera that can be specified from a single string").Out
)

func buildCam(L *l.LState, from int) (*cam, error) {
	var c *cam
	var err error

	stringFn := func(L *l.LState, from int) (*cam, error) {
		var ret *cam
		var serr error
		s := L.CheckString(from)
		switch strings.ToLower(s) {
		case "perspective":
			ret = perspective(65, Aspect(nativeWindow), 0.01, 1000)
		case "orthographic":
			ret = orthographic(-2, 2, 2, -2, 0.01, 100)
		default:
			serr = CameraFromStringError(from)
		}
		return ret, serr
	}

	tableFn := func(L *l.LState, from int) (*cam, error) {
		//t := L.CheckTable(from)
		//spew.Dump(t)
		return nil, nil
	}

	v := L.Get(from)
	switch v.Type() {
	case l.LTString:
		c, err = stringFn(L, from)
	case l.LTTable:
		c, err = tableFn(L, from)
	default:
		L.RaiseError("A string denoting a camera type or table with camera configuration expected. %s is neither", v)
	}
	return c, err
}

func lcamera(L *l.LState) int {
	tag := cameraTag(L)
	cm, err := buildCam(L, 2)
	if err != nil {
		L.RaiseError("error building camera: %s", err)
		return 0
	}
	c := Camera(tag, cm)
	return pushNode(L, c)
}

type cameraMemberFunc func(*l.LState, *l.LUserData, *cameraNode) int

func checkCameraNodeWithUD(L *l.LState, pos int) (*l.LUserData, *cameraNode) {
	ud := L.CheckUserData(pos)
	if n, ok := ud.Value.(*cameraNode); ok {
		return ud, n
	}
	L.ArgError(pos, "camera node expected")
	return nil, nil
}

func cameraMember(fn cameraMemberFunc) l.LGFunction {
	return func(L *l.LState) int {
		if u, n := checkCameraNodeWithUD(L, 1); n != nil {
			return fn(L, u, n)
		}
		return 0
	}
}

func cameraFov(L *l.LState, u *l.LUserData, n *cameraNode) int {
	return 0
}

func cameraAspect(L *l.LState, u *l.LUserData, n *cameraNode) int {
	return 0
}

func cameraNear(L *l.LState, u *l.LUserData, n *cameraNode) int {
	return 0
}

func cameraFar(L *l.LState, u *l.LUserData, n *cameraNode) int {
	return 0
}

func cameraZoom(L *l.LState, u *l.LUserData, n *cameraNode) int {
	return 0
}

var lCameraNodeTable = &lua.Table{
	lCameraNodeClass,
	[]*lua.Table{nodeTable},
	defaultIdxMetaFuncs(),
	map[string]l.LGFunction{},
	map[string]l.LGFunction{
		"fov":    cameraMember(cameraFov),
		"aspect": cameraMember(cameraAspect),
		"near":   cameraMember(cameraNear),
		"far":    cameraMember(cameraFar),
		"zoom":   cameraMember(cameraZoom),
		//view matrix
		//projection matrix
	},
}

// A default perspective camera.
func Perspective() {}

// A default orthographic camera.
func Orthographic() {}
