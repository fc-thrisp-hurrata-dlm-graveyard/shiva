package scene

/*
const lDrawNodeClass = "NDRAW"

type drawNode struct {
	*node
}

func Draw(tag string) Node {
	return &drawNode{
		newNode(tag, func(r *render.Renderer, n Node) {
			//
		}, defaultRemovalFn, defaultReplaceFn, lDrawNodeClass, lNodeClass),
	}
}

func checkDrawNode(L *l.LState) (*l.LUserData, *drawNode) {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*drawNode); ok {
		return ud, v
	}
	L.ArgError(1, "draw node expected")
	return nil, nil
}

type drawPropertyFunc func(*l.LState, *l.LUserData, *drawNode) int

func drawPropertyFn(fn drawPropertyFunc) l.LGFunction {
	return func(L *l.LState) int {
		if ud, n := checkDrawNode(L); n != nil {
			return fn(L, ud, n)
		}
		return 0
	}
}

func ldraw(L *l.LState) int {
	return 0
}

var drawNodeTable = &lua.Table{
	lDrawNodeClass,
	[]*lua.Table{nodeTable},
	nil,
	map[string]l.LGFunction{},
	map[string]l.LGFunction{},
}
*/
