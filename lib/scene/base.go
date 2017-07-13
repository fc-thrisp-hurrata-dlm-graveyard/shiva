package scene

import (
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	"github.com/Laughs-In-Flowers/shiva/lib/render"

	l "github.com/yuin/gopher-lua"
)

const lDummyNodeClass = "NDUMMY"

func Dummy(tag string) Node {
	return newNode(tag, func(r render.Renderer, n Node) {
		for _, nd := range n.Out() {
			nd.Render(r)
		}
	}, defaultRemovalFn, defaultReplaceFn, lDummyNodeClass, lNodeClass)
}

func ldummy(L *l.LState) int {
	tagFn := tagFnFor("dummy", 1)
	d := Dummy(tagFn(L))
	return pushNode(L, d)
}

const lRootNodeClass = "NROOT"

func Root(tag string) Node {
	return newNode(tag, func(r render.Renderer, n Node) {
		for _, nd := range n.Out() {
			nd.Render(r)
		}
	}, defaultRemovalFn, defaultReplaceFn, lRootNodeClass, lNodeClass)
}

func lroot(L *l.LState) int {
	tagFn := tagFnFor("root", 1)
	r := Root(tagFn(L))
	return pushNode(L, r)
}

const lGroupNodeClass = "NGROUP"

func Group(tag string, ns ...Node) Node {
	return newNode(tag, func(r render.Renderer, n Node) {
		for _, nd := range n.Out() {
			nd.Render(r)
		}
	}, defaultRemovalFn, defaultReplaceFn, lGroupNodeClass, lNodeClass)
}

func lgroup(L *l.LState) int {
	var tag string
	var ns []Node

	tagFn := tagFnFor("group", 1)

	tableFn := func(L *l.LState, in []Node) []Node {
		t := L.CheckTable(1)
		t.ForEach(func(k, v l.LValue) {
			if ud, ok := v.(*l.LUserData); ok {
				if vv, ok := ud.Value.(Node); ok {
					in = append(in, vv)
				}
			}
		})
		return in
	}

	fromTop := func(L *l.LState, pos int, in []Node) []Node {
		for i := pos; i <= L.GetTop(); i++ {
			in = append(in, pullNode(L, i))
		}
		return in
	}

	top := L.GetTop()
	switch {
	case top == 1:
		v := L.Get(1)
		switch v.Type() {
		case l.LTString:
			tag = tagFn(L)
		case l.LTTable:
			ns = tableFn(L, ns)
		}
	case top > 1:
		tag = tagFn(L)
		ns = fromTop(L, 2, ns)
	default:
		L.RaiseError("at least 1 argument expected(string tag, nodes list, table of nodes or any combination of nodes and node tables, starting with a string tag)")
	}
	g := Group(tag, ns...)
	return pushNode(L, g)
}

func groupNodeCall(t *lua.Table, k string) l.LGFunction {
	return nodeMember(nodeAppendOut)
}
