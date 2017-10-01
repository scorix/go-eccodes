package native

/*
#include <eccodes.h>
*/
import "C"
import (
	"unsafe"

	"github.com/pkg/errors"
)

func Ccodes_keys_iterator_new(handle Ccodes_handle, flags int, namespace string) Ccodes_keys_iterator {
	cNamespace := C.CString(namespace)
	defer C.free(unsafe.Pointer(cNamespace))

	return unsafe.Pointer(C.codes_keys_iterator_new((*C.codes_handle)(handle), C.ulong(Culong(flags)), cNamespace))
}

func Ccodes_keys_iterator_delete(kiter Ccodes_keys_iterator) error {
	err := C.codes_keys_iterator_delete((*C.codes_keys_iterator)(kiter))
	if err != 0 {
		return errors.New(Cgrib_get_error_message(int(err)))
	}
	return nil
}