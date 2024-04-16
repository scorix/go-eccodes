package codes

import (
	"runtime"

	"github.com/scorix/go-eccodes/debug"
	"github.com/scorix/go-eccodes/native"
)

type Iterator interface {
	Next() (latitude float64, longitude float64, value float64, hasNext bool)
	Previous() (latitude float64, longitude float64, value float64, hasPrevious bool)
	HasNext() bool
	Reset() error

	isOpen() bool
	Close() error
}

type iterator struct {
	it native.Ccodes_iterator
}

func newIterator(handle native.Ccodes_handle) (*iterator, error) {
	it, err := native.Ccodes_grib_iterator_new(handle, 0)
	if err != nil {
		return nil, err
	}

	i := &iterator{it: it}
	runtime.SetFinalizer(i, iteratorFinalizer)

	return i, nil
}

func (it *iterator) Next() (latitude float64, longitude float64, value float64, hasNext bool) {
	return native.Ccodes_grib_iterator_next(it.it)
}

func (it *iterator) Previous() (latitude float64, longitude float64, value float64, hasPrevious bool) {
	return native.Ccodes_grib_iterator_previous(it.it)
}

func (it *iterator) HasNext() bool {
	return native.Ccodes_grib_iterator_has_next(it.it)
}

func (it *iterator) Reset() error {
	return native.Ccodes_grib_iterator_reset(it.it)
}

func (it *iterator) isOpen() bool {
	return it.it != nil
}

func (it *iterator) Close() error {
	defer func() { it.it = nil }()
	return native.Ccodes_grib_iterator_delete(it.it)
}

func iteratorFinalizer(i *iterator) {
	if i.isOpen() {
		debug.MemoryLeakLogger.Print("iterator is not closed")
		i.Close()
	}
}
