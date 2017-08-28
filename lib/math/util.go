package math

import (
	"strconv"
	"strings"

	"github.com/Laughs-In-Flowers/shiva/lib/xrror"
	l "github.com/yuin/gopher-lua"
)

var NoXOfSizeError = xrror.Xrror("No %s of size %d is possible").Out

func defaultToErrMsgs(k string) toErrMsgs {
	return toErrMsgs{
		"%s is not UserData",
		strings.Join([]string{"%s is not", k}, " "),
	}
}

type toErrMsgs [2]string

func establishLimit(k string) int {
	var limit int
	switch k {
	case VEC2, MAT2:
		limit = 2
	case VEC3, MAT3:
		limit = 3
	case VEC4, MAT4:
		limit = 4
	}
	return limit
}

func establishKindOfLimit(k string, l int) string {
	var kind string
	switch k {
	case "vector":
		switch l {
		case 2:
			kind = VEC2
		case 3:
			kind = VEC3
		case 4:
			kind = VEC4
		}
	case "matrice":
		switch l {
		case 2:
			kind = MAT2
		case 3:
			kind = MAT3
		case 4:
			kind = MAT4
		}
	}
	return kind
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TODO: proper LNumber formatting, now its either or bullshit take your pick
// TODO: this is not that cool
func Pf32(L *l.LState, pos int) float32 {
	var num float64
	var err error
	n := L.CheckNumber(pos)
	num, err = strconv.ParseFloat(n.String(), 64)
	if err != nil {
		L.RaiseError(err.Error())
	}
	return float32(num)
}

type MxPos struct {
	Correspondence, Row, Column int
}

var (
	// vector x,y,z to matrice raw[12,13,14]
	TranslateMxPos = []MxPos{
		{0, 0, 3},
		{1, 1, 3},
		{2, 2, 3},
	}
)

// TODO: this is a fucking piece of uncooperative shit
//type fmts struct {
//value     float64
//width     int
//precision int
//}

//func newFmts() *fmts {
//	return &fmts{0, 16, 16}
//}

//func (f *fmts) Write(b []byte) (n int, err error) {
//spew.Dump(b)
//buf := bytes.NewReader(b) //new(bytes.Buffer)
//binary.Write(buf, binary.LittleEndian, b)
//spew.Dump(buf)
//binary.Read(buf, binary.LittleEndian, &f.value)
//fmt.Print(f.value)
//spew.Dump(f)
//spew.Dump(fmt.Sprintf("%G", f.value))
//bits := binary.LittleEndian.Uint64(b)
//float := math.Float64frombits(bits)
//spew.Dump(float)
//return 0, nil
//}

//func (f *fmts) Width() (wid int, ok bool) {
//	return f.width, true
//}

//func (f *fmts) Precision() (prec int, ok bool) {
//	return f.precision, true
//}

//func (f *fmts) Flag(c int) bool {
//	return false
//}
