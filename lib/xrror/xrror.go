package xrror

import (
	"fmt"
	"os"
)

func Basic(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", e)
		os.Exit(-1)
	}
}

type xrror struct {
	base string
	vals []interface{}
}

func (x *xrror) Error() string {
	return fmt.Sprintf("%s", fmt.Sprintf(x.base, x.vals...))
}

func (x *xrror) Out(vals ...interface{}) *xrror {
	x.vals = vals
	return x
}

func Xrror(base string) *xrror {
	return &xrror{base: base}
}
