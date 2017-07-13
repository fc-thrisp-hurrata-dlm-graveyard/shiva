package lua

import (
	l "github.com/yuin/gopher-lua"
)

type RegisterWith func(Module) error

type GetModule func(*Lua) Module

type Module interface {
	Tag() string
	AddLGFunc(string, l.LGFunction)
	Loader(*l.LState) int
	AddMT(...MLoadMT)
	LoadMT(*l.LState)
	Register(*l.LState, *Table)
	AddEmb(...MLoadEmb)
	LoadEmb(*Lua) error
}

type MLoader func(*l.LState) int

type MLoadMT func(*l.LState, Module)

type MLoadEmb func(*Lua) error

type module struct {
	name       string
	loaderFn   MLoader
	loadFns    map[string]l.LGFunction
	loadMts    []MLoadMT
	metatables *Tables
	loadEmb    []MLoadEmb
}

func NewModule(name string, mlr MLoader) Module {
	return &module{
		name,
		mlr,
		make(map[string]l.LGFunction),
		make([]MLoadMT, 0),
		NewTables(),
		make([]MLoadEmb, 0),
	}
}

func (m *module) Tag() string {
	return m.name
}

func (m *module) Loader(L *l.LState) int {
	return m.loaderFn(L)
}

func (m *module) AddLGFunc(name string, fn l.LGFunction) {
	m.loadFns[name] = fn
}

func (m *module) AddMT(fns ...MLoadMT) {
	m.loadMts = append(m.loadMts, fns...)
}

func (m *module) LoadMT(L *l.LState) {
	for _, mt := range m.loadMts {
		mt(L, m)
	}
}

func (m *module) Register(L *l.LState, t *Table) {
	m.metatables.Register(L, t)
}

func (m *module) AddEmb(fns ...MLoadEmb) {
	m.loadEmb = append(m.loadEmb, fns...)
}

func (m *module) LoadEmb(L *Lua) error {
	var err error
	for _, fn := range m.loadEmb {
		err = fn(L)
		if err != nil {
			return err
		}
	}
	return err
}

var defaultModules = []GetModule{
	func(*Lua) Module { return NewModule(l.LoadLibName, l.OpenPackage) },
	func(*Lua) Module { return NewModule(l.BaseLibName, l.OpenBase) },
	func(*Lua) Module { return NewModule("string", l.OpenString) },
	func(*Lua) Module { return NewModule("table", l.OpenTable) },
	func(*Lua) Module { return NewModule("os", l.OpenOs) },
	func(*Lua) Module { return NewModule("io", l.OpenIo) },
	func(*Lua) Module { return NewModule("debug", l.OpenDebug) },
	shvGetM,
}
