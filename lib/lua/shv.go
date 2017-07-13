package lua

import (
	"bytes"
	"fmt"

	l "github.com/yuin/gopher-lua"
)

const ShvModuleName = "shv"

var (
	ShvModule  Module = mkShvModule()
	lShvModule *l.LTable
	shvGetM    = func(*Lua) Module { return ShvModule }
)

func mkShvModule() Module {
	m := &module{
		name:       ShvModuleName,
		loadFns:    make(map[string]l.LGFunction),
		loadMts:    make([]MLoadMT, 0),
		metatables: NewTables(),
		loadEmb:    make([]MLoadEmb, 0),
	}
	m.loaderFn = moduleLoader(m)
	m.AddEmb(shvLoadEmb)
	return m
}

func moduleLoader(m *module) func(L *l.LState) int {
	return func(L *l.LState) int {
		mod := L.RegisterModule(m.name, m.loadFns).(*l.LTable)
		L.SetField(mod, "name", l.LString(m.name))
		for _, mt := range m.loadMts {
			mt(L, m)
		}
		lShvModule = mod
		L.Push(mod)
		return 1
	}
}

func shvLoadEmb(L *Lua) error {
	a := ResourceFS
	files := a.AssetNames()
	for _, f := range files {
		var b []byte
		var err error
		b, err = a.Asset(f)
		if err != nil {
			return err
		}
		bf := bytes.NewBuffer(b)
		n := fmt.Sprintf("embedded-%s", f)
		var rfn *l.LFunction
		rfn, err = L.Load(bf, n)
		if err != nil {
			return err
		}
		L.Push(rfn)
		L.CallWithTraceback(0, 0)
	}
	return nil
}

func (L *Lua) CallShv(fn string, nargs, nresults int) error {
	lfn := L.GetField(lShvModule, fn)
	L.Remove(-2)
	if nargs > 0 {
		L.Insert(lfn, -nargs-1)
	}
	return L.CallWithTraceback(nargs, nresults)
}
