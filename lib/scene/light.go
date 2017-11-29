package scene

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/Laughs-In-Flowers/shiva/lib/math"
	"github.com/Laughs-In-Flowers/shiva/lib/render"

	l "github.com/yuin/gopher-lua"
)

type LightKind int

const (
	NOT_A_LIGHT LightKind = iota
	AMBIENT
	DIRECTIONAL
	POINT
	SPOT
)

var lightSet = []LightKind{AMBIENT, DIRECTIONAL, POINT, SPOT}

type indexLights map[LightKind]int

func (il indexLights) current(k LightKind) int {
	if v, exists := il[k]; exists {
		il[k] = v + 1
		return il[k]
	}
	return 0
}

func newIndexLights() indexLights {
	ret := make(map[LightKind]int)
	for _, k := range lightSet {
		ret[k] = 0
	}
	return ret
}

var LightIndex indexLights

func init() {
	LightIndex = newIndexLights()
}

type light struct {
	*node
	idx       int
	intensity float32
	color     math.Color
	u         graphics.Uniform
}

func (t *light) Initialize(k LightKind) {
	t.SetColor(t.color)
	t.SetIntensity(t.intensity)
	t.SetIdx(LightIndex.current(k))
}

func (t *light) Idx() int {
	return t.idx
}

func (t *light) SetIdx(i int) {
	t.idx = i
}

func (t *light) Color() math.Color {
	return t.color
}

func (t *light) SetColor(c math.Color) {
	t.color = c
	t.postChange()
}

func (t *light) Intensity() float32 {
	return t.intensity
}

func (t *light) SetIntensity(v float32) {
	t.intensity = v
	t.postChange()
}

func (t *light) postChange() {
	t.color.MulScalar(t.intensity)
	t.u.Update(t.color.Raw()...)
}

func (t *light) Provide(p graphics.Provider) {
	t.u.TransferIdx(p, t.idx)
}

const lLightNodeClass = "NLIGHT"

var lLightNodeTable = &lua.Table{
	lAmbientLightNodeClass,
	[]*lua.Table{nodeTable},
	nil, nil, nil,
}

const lAmbientLightNodeClass = "NLAMBIENT"

func Ambient(tag string, intensity float32, color math.Color) *light {
	u := graphics.Uniform3f("AmbientLightColor")
	lg := &light{nil, 0, intensity, color, u}
	n := newNode(tag, func(r render.Renderer, n Node) {
		lg.Provide(r)
	}, defaultRemovalFn, defaultReplaceFn, lAmbientLightNodeClass, lLightNodeClass, lNodeClass)
	n.lightKind = AMBIENT
	lg.node = n
	lg.Initialize(n.lightKind)
	return lg
}

func lambient(L *l.LState) int {
	//tag intensity color
	return 0
}

var lAmbientLightNodeTable = &lua.Table{
	lAmbientLightNodeClass,
	[]*lua.Table{lLightNodeTable},
	nil, nil, nil,
}

const lDirectionalLightNodeClass = "NLDIRECTIONAL"

func Directional(tag string) *light {
	return nil
}

func ldirectional(L *l.LState) int {
	return 0
}

var lDirectionalLightNodeTable = &lua.Table{
	lDirectionalLightNodeClass,
	[]*lua.Table{lLightNodeTable},
	nil, nil, nil,
}

const lPointLightNodeClass = "NLPOINT"

func Point(tag string) *light {
	return nil
}

func lpoint(L *l.LState) int {
	return 0
}

var lPointLightNodeTable = &lua.Table{
	lPointLightNodeClass,
	[]*lua.Table{lLightNodeTable},
	nil, nil, nil,
}

const lSpotLightNodeClass = "NLSPOT"

func Spot(tag string) *light {
	return nil
}

func lspot(L *l.LState) int {
	return 0
}

var lSpotLightNodeTable = &lua.Table{
	lSpotLightNodeClass,
	[]*lua.Table{lLightNodeTable},
	nil, nil, nil,
}
