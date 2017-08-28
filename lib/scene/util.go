package scene

import (
	l "github.com/yuin/gopher-lua"
)

type TagFunc func(L *l.LState) string

func tagFnFor(s string, pos int) TagFunc {
	return func(L *l.LState) string {
		var ret = s
		rtag := L.ToString(pos)
		if rtag != "" {
			ret = rtag
		}
		return ret
	}
}
