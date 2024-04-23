//go:build linux && (amd64 || arm64)
// +build linux
// +build amd64 arm64

package native

/*
#include <eccodes.h>
#cgo CFLAGS: -I /usr/share/eccodes/include
#cgo LDFLAGS: -L /usr/share/eccodes/lib -leccodes
*/
import "C"

type (
	Cint    = int32
	Clong   = int64
	Culong  = uint64
	Cdouble = float64
	CsizeT  = int64
)
