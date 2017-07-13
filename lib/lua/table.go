package lua

import l "github.com/yuin/gopher-lua"

type MetaLGFunc func(*Table, string) l.LGFunction

type LMetaFunc struct {
	Key   string
	Value MetaLGFunc
}

type Table struct {
	Name     string
	Parent   []*Table
	Meta     []*LMetaFunc
	Property map[string]l.LGFunction
	Method   map[string]l.LGFunction
}

func NewTable(
	name string,
	parent []*Table,
	meta []*LMetaFunc,
	prop, meth map[string]l.LGFunction,
) *Table {
	return &Table{
		name,
		parent,
		meta,
		prop,
		meth,
	}
}

func (t *Table) accessible(to, key string) (l.LGFunction, bool) {
	switch to {
	case "property":
		if a, exists := t.Property[key]; exists {
			return a, true
		}
		for _, p := range t.Parent {
			if pa, exists := p.accessible("property", key); exists {
				return pa, true
			}
		}
	case "method":
		if a, exists := t.Method[key]; exists {
			return a, true
		}
		for _, p := range t.Parent {
			if pa, exists := p.accessible("method", key); exists {
				return pa, true
			}
		}
	}
	return nil, false
}

func (t *Table) Execute(L *l.LState) int {
	req := L.Get(2)
	k := req.String()
	if property, exists := t.accessible("property", k); exists {
		return property(L)
	}
	if method, exists := t.accessible("method", k); exists {
		nfn := L.NewFunction(method)
		L.Push(nfn)
		return 1
	}
	L.RaiseError("%s is not a member of %s or any table %s is descended from", k, t.Name, t.Name)
	return 0
}

type Tables struct {
	has map[string]*Table
}

func NewTables() *Tables {
	return &Tables{
		make(map[string]*Table),
	}
}

func (t *Tables) Get(k string) *Table {
	if tbl, ok := t.has[k]; ok {
		return tbl
	}
	return nil
}

func (t *Tables) Register(L *l.LState, tbl *Table) {
	t.Set(tbl)
	mt := L.NewTypeMetatable(tbl.Name)
	for _, v := range tbl.Meta {
		key := v.Key
		fn := v.Value
		mkTableFunc(L, mt, key, fn(tbl, key))
	}
}

func mkTableFunc(L *l.LState, t *l.LTable, k string, fn l.LGFunction) {
	L.SetField(t, k, L.NewClosure(fn))
}

func (t *Tables) Set(tbl *Table) {
	if _, exists := t.has[tbl.Name]; !exists {
		t.has[tbl.Name] = tbl
	}
}

func defaultIdxFunc(t *Table, k string) l.LGFunction {
	return func(L *l.LState) int {
		return t.Execute(L)
	}
}

func DefaultIdx(tag string) *LMetaFunc {
	return &LMetaFunc{
		tag,
		defaultIdxFunc,
	}
}

func immutableNewIndexFunc() MetaLGFunc {
	return func(t *Table, k string) l.LGFunction {
		return func(L *l.LState) int {
			L.RaiseError("immutable %s: attempt to add new index", t.Name)
			return 0
		}
	}
}

func ImmutableNewIdx() *LMetaFunc {
	return &LMetaFunc{
		"__newindex",
		immutableNewIndexFunc(),
	}
}
