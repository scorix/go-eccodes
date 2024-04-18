//go:build darwin
// +build darwin

package native

/*
#include <eccodes.h>
#cgo CFLAGS: -I /opt/homebrew/include
#cgo LDFLAGS: -L /opt/homebrew/lib -leccodes
*/
import "C"

type (
	Cint    = int32
	Clong   = int64
	Culong  = uint64
	Cdouble = float64
	CsizeT  = int64
)
