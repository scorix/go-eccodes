//go:build darwin
// +build darwin

package native

/*
#include <eccodes.h>
#cgo CFLAGS: -I /opt/homebrew/include
#cgo LDFLAGS: -L /opt/homebrew/lib -leccodes
*/
import "C"

type Cint = int32
type Clong = int64
type Culong = uint64
type Cdouble = float64
type CsizeT = int64
