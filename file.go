package codes

import (
	"errors"
	"fmt"
	"io"
	"runtime"

	"github.com/scorix/go-eccodes/debug"
	cio "github.com/scorix/go-eccodes/io"
	"github.com/scorix/go-eccodes/native"
)

type Reader interface {
	Next() (Message, error)
}

type Writer interface{}

type File interface {
	Reader
	Writer
	Close()
}

type file struct {
	file cio.File
}

type fileIndexed struct {
	index native.Ccodes_index
}

// https://confluence.ecmwf.int/display/ECC/grib_copy
var emptyFilter = map[string]any{}

func OpenFile(f cio.File) (File, error) {
	return &file{file: f}, nil
}

func OpenFileByPathWithFilter(path string, filter map[string]any) (File, error) {
	if filter == nil {
		filter = emptyFilter
	}

	var k string
	for key, value := range filter {
		if len(k) > 0 {
			k += ","
		}
		k += key
		if value != nil {
			switch value.(type) {
			case int64, int:
				k += ":l"
			case float64, float32:
				k += ":d"
			case string:
				k += ":s"
			}
		}
	}

	i, err := native.Ccodes_index_new_from_file(native.DefaultContext, path, k)
	if err != nil {
		return nil, fmt.Errorf("failed to create filtered index: %w", err)
	}

	for key, value := range filter {
		switch value := value.(type) {
		case int64:
			if err := native.Ccodes_index_select_long(i, key, value); err != nil {
				native.Ccodes_index_delete(i)

				return nil, fmt.Errorf("failed to set filter condition '%s'=%d: %w", key, value, err)
			}
		case int:
			if err := native.Ccodes_index_select_long(i, key, int64(value)); err != nil {
				native.Ccodes_index_delete(i)

				return nil, fmt.Errorf("failed to set filter condition '%s'=%d: %w", key, value, err)
			}
		case float64:
			if err := native.Ccodes_index_select_double(i, key, value); err != nil {
				native.Ccodes_index_delete(i)

				return nil, fmt.Errorf("failed to set filter condition '%s'=%f: %w", key, value, err)
			}
		case float32:
			if err := native.Ccodes_index_select_double(i, key, float64(value)); err != nil {
				native.Ccodes_index_delete(i)

				return nil, fmt.Errorf("failed to set filter condition '%s'=%f: %w", key, value, err)
			}
		case string:
			if err := native.Ccodes_index_select_string(i, key, value); err != nil {
				native.Ccodes_index_delete(i)

				return nil, fmt.Errorf("failed to set filter condition '%s'=%q: %w", key, value, err)
			}
		}
	}

	file := &fileIndexed{index: i}
	runtime.SetFinalizer(file, fileIndexedFinalizer)

	return file, nil
}

func (f *file) Next() (Message, error) {
	handle, err := native.Ccodes_handle_new_from_file(native.DefaultContext, f.file.Native(), native.ProductAny)
	if errors.Is(err, io.EOF) {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed create new handle from file: %w", err)
	}

	return newMessage(handle), nil
}

func (f *file) Close() {
	f.file = nil
}

func (f *fileIndexed) isOpen() bool {
	return f.index != nil
}

func (f *fileIndexed) Next() (Message, error) {
	handle, err := native.Ccodes_handle_new_from_index(f.index)
	if errors.Is(err, io.EOF) {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create handle from index: %w", err)
	}

	return newMessage(handle), nil
}

func (f *fileIndexed) Close() {
	if f.isOpen() {
		defer func() { f.index = nil }()
		native.Ccodes_index_delete(f.index)
	}
}

func fileIndexedFinalizer(f *fileIndexed) {
	if f.isOpen() {
		debug.MemoryLeakLogger.Print("file is not closed")
		f.Close()
	}
}
