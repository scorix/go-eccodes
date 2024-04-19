//go:build linux && (amd64 || arm64)
// +build linux
// +build amd64 arm64

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
