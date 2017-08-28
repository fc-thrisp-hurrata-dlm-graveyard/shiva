package scene

import (
	"reflect"

	"github.com/Laughs-In-Flowers/shiva/lib/ecs"
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/Laughs-In-Flowers/shiva/lib/render"
	"github.com/Laughs-In-Flowers/shiva/lib/xrror"

	l "github.com/yuin/gopher-lua"
)

const lNodeClass = "NODE"

type Node interface {
	Tagger
	Recurser
	Renderable
	Actor
	Grapher
	Comparer
	Counter
	Terminater
	//Removal
}

type Tagger interface {
	ecs.Entity
	Tag() string
	Retag(string)
	LClass() string
}

type RenderFunc func(render.Renderer)

type innerRenderFunc func(render.Renderer, Node)

type Recurser interface {
	Hidden() bool
	SetHidden(bool)
	RecursionLimit() int8
	SetRecursionLimit(int8)
}

type recurser struct {
	recursed, recurse int8
	hidden            bool
}

func newRecurser() *recurser {
	return &recurser{}
}

func (r *recurser) rExecute(fn func()) {
	if r.recurse == 0 || r.recursed <= r.recurse {
		if !r.hidden {
			fn()
			r.recursed++
		}
	}
}

func (r *recurser) Hidden() bool {
	return r.hidden
}

func (r *recurser) SetHidden(hidden bool) {
	r.hidden = hidden
}

func (r *recurser) RecursionLimit() int8 {
	return r.recurse
}

func (r *recurser) SetRecursionLimit(rl int8) {
	r.recurse = rl
}

type Renderable interface {
	Render(render.Renderer)
}

type Actor interface {
	Paused() bool
	SetPaused(bool)
	Act(*l.LState) error
}

type actor struct {
	acted, act int8
	paused     bool
	actions    []l.LGFunction
}

func newActor() *actor {
	return &actor{
		actions: make([]l.LGFunction, 0),
	}
}

func (a *actor) Paused() bool {
	return a.paused
}

func (a *actor) SetPaused(paused bool) {
	a.paused = paused
}

func (a *actor) Act(L *l.LState) error {
	return nil
}

type NIter func(Node)

type Grapher interface {
	In() []Node
	Out() []Node
	Append(RelationDir, ...Node) error
	Prepend(RelationDir, ...Node) error
	Remove()
	Replace(Node)
	Iter(RelationDir, NIter)
}

type Comparer interface {
	Equal(Node) bool
}

type Counter interface {
	Count() int
}

type Terminater interface {
	Terminal() bool
	SetTerminal(bool)
}

type Lighter interface {
	Light() LightKind
}

type removalFunc func(*node) error

type replaceFunc func(*node, *node) error

type node struct {
	ecs.Entity
	*recurser
	*actor
	tag       string
	lclass    []string
	rfn       innerRenderFunc
	rmfn      removalFunc
	rpfn      replaceFunc
	in        []Node
	out       []Node
	lightKind LightKind
	terminal  bool
}

func newNode(
	tag string,
	fn innerRenderFunc,
	rmfn removalFunc,
	rpfn replaceFunc,
	lclass ...string,
) *node {
	return &node{
		Entity:   ecs.NewEntity(),
		recurser: newRecurser(),
		actor:    newActor(),
		tag:      tag,
		lclass:   lclass,
		rfn:      fn,
		rmfn:     rmfn,
		rpfn:     rpfn,
		in:       make([]Node, 0),
		out:      make([]Node, 0),
	}
}

func (n *node) Tag() string {
	return n.tag
}

func (n *node) Retag(nt string) {
	n.tag = nt
}

func (n *node) LClass() string {
	if len(n.lclass) >= 1 {
		return n.lclass[0]
	}
	return lNodeClass
}

func (n *node) Render(r render.Renderer) {
	n.rExecute(func() {
		n.rfn(r, n)
	})
	if !n.terminal {
		for _, nd := range n.Out() {
			nd.Render(r)
		}
	}
}

func (n *node) In() []Node {
	return n.in
}

func (n *node) Out() []Node {
	return n.out
}

type RelationDir int

const (
	NIN RelationDir = iota
	NOUT
	UNKNOWNDIR
)

func (r RelationDir) String() string {
	switch r {
	case NIN:
		return "IN"
	case NOUT:
		return "OUT"
	}
	return "UNKNOWN"
}

func stringRelationDir(s string) RelationDir {
	switch s {
	case "IN", "in":
		return NIN
	case "OUT", "out":
		return NOUT
	}
	return UNKNOWNDIR
}

var NodeRelationError = xrror.Xrror("%s is not a direction for node relation %s").Out

func (n *node) Append(dir RelationDir, ns ...Node) error {
	switch dir {
	case NIN:
		n.in = append(n.in, ns...)
		return nil
	case NOUT:
		for _, v := range ns {
			v.Append(NIN, n)
		}
		n.out = append(n.out, ns...)
		return nil
	}
	return NodeRelationError(dir.String(), "append")
}

func (n *node) Prepend(dir RelationDir, ns ...Node) error {
	switch dir {
	case NIN:
		nin := append(ns, n.in...)
		n.in = nin
		return nil
	case NOUT:
		for _, v := range ns {
			v.Append(NIN, n)
		}
		nout := append(ns, n.out...)
		n.out = nout
		return nil
	}
	return NodeRelationError(dir.String(), "prepend")
}

func defaultRemovalFn(n *node) error {
	eraseOut := func(p, c *node) {
		for idx, cn := range p.out {
			if cn.Equal(c) {
				copy(p.out[idx:], p.out[idx+1:])
				p.out[len(p.out)-1] = nil
				p.out = p.out[:len(p.out)-1]
			}
		}
	}
	for _, parent := range n.in {
		if pn, ok := parent.(*node); ok {
			eraseOut(pn, n)
		}
	}
	return nil
}

func (n *node) Remove() {
	n.rmfn(n)
}

func defaultReplaceFn(n, with *node) error {
	replaceOut := func(p, c, w *node) {
		for idx, cn := range p.out {
			if cn.Equal(c) {
				p.out[idx] = w
			}
		}
	}
	for _, parent := range n.in {
		if pn, ok := parent.(*node); ok {
			replaceOut(pn, n, with)
		}
	}
	return nil
}

func (n *node) Replace(with Node) {
	if w, ok := with.(*node); ok {
		n.rpfn(n, w)
	}
}

func (n *node) Iter(dir RelationDir, fn NIter) {
	fn(n)
	switch dir {
	case NIN:
		for _, nn := range n.in {
			nn.Iter(NIN, fn)
		}
	case NOUT:
		for _, nn := range n.out {
			nn.Iter(NOUT, fn)
		}
	}
}

func (n *node) Equal(o Node) bool {
	return reflect.DeepEqual(n, o)
}

func (n *node) Count() int {
	count := 0
	n.Iter(NOUT, func(n Node) {
		count = count + 1
	})
	return count
}

func (n *node) Terminal() bool {
	return n.terminal
}

func (n *node) SetTerminal(as bool) {
	n.terminal = as
}

func pushNode(L *l.LState, n Node) int {
	fn := func(u *l.LUserData) {
		u.Value = n
	}
	lua.PushNewUserData(L, fn, n.LClass())
	return 1
}

func pullNodeUD(L *l.LState, pos int) (*l.LUserData, Node) {
	ud := L.CheckUserData(pos)
	if v, ok := ud.Value.(Node); ok {
		return ud, v
	}
	L.ArgError(pos, "node expected")
	return nil, nil
}

func pullNode(L *l.LState, pos int) Node {
	_, n := pullNodeUD(L, pos)
	if n != nil {
		return n
	}
	return nil
}

func nodeChain(t *lua.Table, k string) l.LGFunction {
	return func(L *l.LState) int {
		parent := checkNode(L, 1)
		child := checkNode(L, 2)
		parent.Append(NOUT, child)
		return pushNode(L, parent)
	}
}

func isBaseNode(n Node) (*node, bool) {
	if nd, ok := n.(*node); ok {
		return nd, true
	}
	return nil, false
}

func isNode(v l.LValue) (Node, bool) {
	if ud, ok := v.(*l.LUserData); ok {
		if n, ok := ud.Value.(Node); ok {
			return n, true
		}
	}
	return nil, false
}

func checkNode(L *l.LState, pos int) Node {
	_, n := checkNodeWithUD(L, pos)
	return n
}

func checkNodeWithUD(L *l.LState, pos int) (*l.LUserData, Node) {
	ud := L.CheckUserData(pos)
	if n, ok := ud.Value.(Node); ok {
		return ud, n
	}
	L.ArgError(pos, "node expected")
	return nil, nil
}

type nodeMemberFunc func(*l.LState, *l.LUserData, Node) int

func nodeMember(fn nodeMemberFunc) l.LGFunction {
	return func(L *l.LState) int {
		if u, n := checkNodeWithUD(L, 1); n != nil {
			return fn(L, u, n)
		}
		return 0
	}
}

func nodeProperty(get, set nodeMemberFunc) l.LGFunction {
	return lua.NewProperty(nodeMember(get), nodeMember(set))
}

func getHidden(L *l.LState, u *l.LUserData, n Node) int {
	L.Push(l.LBool(n.Hidden()))
	return 1
}

func setHidden(L *l.LState, u *l.LUserData, n Node) int {
	hidden := L.CheckBool(3)
	n.SetHidden(hidden)
	return 0
}

func getPaused(L *l.LState, u *l.LUserData, n Node) int {
	L.Push(l.LBool(n.Paused()))
	return 1
}

func setPaused(L *l.LState, u *l.LUserData, n Node) int {
	paused := L.CheckBool(3)
	n.SetPaused(paused)
	return 0
}

func getTag(L *l.LState, u *l.LUserData, n Node) int {
	L.Push(l.LString(n.Tag()))
	return 1
}

func setTag(L *l.LState, u *l.LUserData, n Node) int {
	tag := L.Get(3)
	n.Retag(tag.String())
	return 0
}

func getRecursionLimit(L *l.LState, u *l.LUserData, n Node) int {
	L.Push(l.LNumber(n.RecursionLimit()))
	return 1
}

func setRecursionLimit(L *l.LState, u *l.LUserData, n Node) int {
	lim := L.CheckInt(3)
	n.SetRecursionLimit(int8(lim))
	return 0
}

func nodeAdd(L *l.LState, afn func(RelationDir, ...Node) error, dir RelationDir) int {
	var add []Node
	ta := L.GetTop()
	for i := 2; i <= ta; i++ {
		vt := L.Get(i)
		switch vt.Type() {
		case l.LTTable:
			tbl := L.CheckTable(i)
			tbl.ForEach(func(k, v l.LValue) {
				if nn, ok := isNode(v); ok {
					add = append(add, nn)
				}
			})
		case l.LTUserData:
			nn := checkNode(L, i)
			add = append(add, nn)
		}
	}
	err := afn(dir, add...)
	if err != nil {
		L.RaiseError(err.Error())
	}
	return 0
}

func nodePrependOut(L *l.LState, u *l.LUserData, n Node) int {
	return nodeAdd(L, n.Prepend, NOUT)
}

func nodeAppendOut(L *l.LState, u *l.LUserData, n Node) int {
	return nodeAdd(L, n.Append, NOUT)
}

func nodeRemove(L *l.LState, u *l.LUserData, n Node) int {
	var ret Node

	niterS := func(s string) func(Node) {
		return func(nn Node) {
			if s == nn.Tag() {
				ret = nn
				nn.Remove()
			}
		}
	}
	niterN := func(cn Node) func(Node) {
		return func(nnn Node) {
			if nnn.Equal(cn) {
				ret = nnn
				nnn.Remove()
			}
		}
	}

	for i := 2; i <= L.GetTop(); i++ {
		v := L.Get(i)
		switch v.Type() {
		case l.LTString:
			s := v.String()
			n.Iter(NOUT, niterS(s))
		case l.LTUserData:
			cn := checkNode(L, i)
			n.Iter(NOUT, niterN(cn))
		default:
			L.ArgError(i, "node or node tag expected")
		}
	}

	if ret != nil {
		pushNode(L, ret)
		return 1
	}

	return 0
}

func nodeReplace(L *l.LState, u *l.LUserData, n Node) int {
	niterN := func(replacement Node) func(Node) {
		return func(nn Node) {
			if nn.Tag() == replacement.Tag() {
				// replace
			}
		}
	}
	for i := 2; i <= L.GetTop(); i++ {
		v := L.Get(i)
		switch v.Type() {
		case l.LTUserData:
			r := checkNode(L, i)
			n.Iter(NOUT, niterN(r))
		default:
			L.ArgError(i, "node expected")
		}
	}
	return 0
}

func nodeAction(L *l.LState, u *l.LUserData, n Node) int {
	// set node action
	return 0
}

func defaultIdxMetaFuncs() []*lua.LMetaFunc {
	return []*lua.LMetaFunc{
		lua.DefaultIdx("__index"),
		lua.DefaultIdx("__newindex"),
		{"__pow", nodeChain},
	}
}

func defaultNodeTable(class string, mfn ...*lua.LMetaFunc) *lua.Table {
	lmf := append(defaultIdxMetaFuncs(), mfn...)
	return &lua.Table{
		class,
		[]*lua.Table{nodeTable},
		lmf,
		nil,
		nil,
	}
}

var nodeTable = &lua.Table{
	lNodeClass,
	nil,
	nil,
	map[string]l.LGFunction{
		"hidden":          nodeProperty(getHidden, setHidden),
		"paused":          nodeProperty(getPaused, setPaused),
		"tag":             nodeProperty(getTag, setTag),
		"recursion_limit": nodeProperty(getRecursionLimit, setRecursionLimit),
	},
	map[string]l.LGFunction{
		"prepend": nodeMember(nodePrependOut),
		"append":  nodeMember(nodeAppendOut),
		"remove":  nodeMember(nodeRemove),
		"replace": nodeMember(nodeReplace),
		"action":  nodeMember(nodeAction),
	},
}

// TODO: tbd nodes
//const lLookatNodeClass = "NLOOKAT"

//func Lookat() Node {
//eye
//center
//up
//return nil
//}

//func llookat(L *l.LState) int {
//return 0
//}

//var lookatNodeTable = &Table{
//lLookatNodeClass,
//nil, nil, nil, nil,
//}

//const lBillboardNodeClass = "NBILLBOARD"

//func Billboard() Node {
//preserve_uniform_scaling bool
//return nil
//}

//func lbillboard(L *l.LState) int {
//return 0
//}

//var billboardNodeTable = &Table{
//lBillboardNodeClass,
//nil, nil, nil, nil,
//}

//blend
//const lBlendNodeClass = "NBLEND"

//color_mask
//const lColorMaskClass = "NCOLORMASK"

//cull_face
//const lCullFaceClass = "NCULLFACE"

//depth_test
//const lDepthTestClass = "NDEPTHTEST"

//viewport
//const lViewportClass = "NVIEWPORT"

//postprocess
//const lPostProcess = "NPOSTPROCESS"
