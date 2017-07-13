package scene

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/lua"
	l "github.com/yuin/gopher-lua"
)

const lProgramNodeClass = "NPROGRAM"

func Program() Node {
	return nil
}

type programNode struct {
	*graphics.Program
}

func lprogram(L *l.LState) int {
	return 0
}

var programNodeTable = &lua.Table{
	lProgramNodeClass,
	nil, nil, nil, nil,
}

const lBindNodeClass = "NBIND"

func Bind() Node {
	return nil
}

type bindNode struct{}

func lbind(L *l.LState) int {
	return 0
}

var bindNodeTable = &lua.Table{
	lBindNodeClass,
	nil, nil, nil, nil,
}
