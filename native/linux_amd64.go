//go:build linux && amd64
// +build linux,amd64

package native

/*
#include <eccodes.h>
#
#cgo LDFLAGS: -leccodes
*/
import "C"

type Cint = int32
type Clong = int64
type Culong = uint64
type Cdouble = float64
type CsizeT = int64
