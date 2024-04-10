package native

/*
#include <eccodes.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

func Ccodes_grib_nearest_new(handle Ccodes_handle) (Ccodes_nearest, error) {
	var err Cint
	cError := (*C.int)(unsafe.Pointer(&err))

	n := C.codes_grib_nearest_new((*C.codes_handle)(handle), cError)
	if err != 0 {
		return nil, errors.New(Cgrib_get_error_message(int(err)))
	}

	return unsafe.Pointer(n), nil
}

func Ccodes_grib_nearest_delete(nearest Ccodes_nearest) error {
	err := C.codes_grib_nearest_delete((*C.codes_nearest)(nearest))
	if err != 0 {
		return errors.New(Cgrib_get_error_message(int(err)))
	}
	return nil
}
