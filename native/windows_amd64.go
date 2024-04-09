//go:build windows && amd64
// +build windows,amd64

package native

/*
#include <eccodes.h>
#
#cgo LDFLAGS: -leccodes -lpng -laec -ljasper -lopenjp2 -lz
*/
import "C"

type Cint = int32
type Clong = int32
type Culong = uint32
type Cdouble = float64
type CsizeT = int64
