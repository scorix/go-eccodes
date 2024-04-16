package native

/*
#include <eccodes.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

func Ccodes_grib_iterator_new(handle Ccodes_handle, flags uint64) (Ccodes_iterator, error) {
	var err Cint
	cError := (*C.int)(unsafe.Pointer(&err))

	i := C.codes_grib_iterator_new((*C.codes_handle)(handle), (C.ulong)(Culong(flags)), cError)
	if err != 0 {
		return nil, errors.New(Cgrib_get_error_message(int(err)))
	}

	return unsafe.Pointer(i), nil
}

func Ccodes_grib_iterator_next(iterator Ccodes_iterator) (latitude float64, longitude float64, value float64, hasNext bool) {
	ret := C.codes_grib_iterator_next((*C.codes_iterator)(iterator), (*C.double)(&latitude), (*C.double)(&longitude), (*C.double)(&value))
	hasNext = ret > 0

	return
}

func Ccodes_grib_iterator_previous(iterator Ccodes_iterator) (latitude float64, longitude float64, value float64, hasPrevious bool) {
	ret := C.codes_grib_iterator_previous((*C.codes_iterator)(iterator), (*C.double)(&latitude), (*C.double)(&longitude), (*C.double)(&value))
	hasPrevious = ret > 0

	return
}

func Ccodes_grib_iterator_has_next(iterator Ccodes_iterator) bool {
	ret := C.codes_grib_iterator_has_next((*C.codes_iterator)(iterator))

	return ret == 1
}

func Ccodes_grib_iterator_reset(iterator Ccodes_iterator) error {
	err := C.codes_grib_iterator_reset((*C.codes_iterator)(iterator))
	if err != 0 {
		return errors.New(Cgrib_get_error_message(int(err)))
	}

	return nil
}

func Ccodes_grib_iterator_delete(iterator Ccodes_iterator) error {
	err := C.codes_grib_iterator_delete((*C.codes_iterator)(iterator))
	if err != 0 {
		return errors.New(Cgrib_get_error_message(int(err)))
	}

	return nil
}
