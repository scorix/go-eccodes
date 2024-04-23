package codes

import (
	"github.com/scorix/go-eccodes/debug"
	"github.com/scorix/go-eccodes/native"
)

type Handle struct {
	handle native.Ccodes_handle
}

func (h *Handle) Message() Message {
	return newMessage(h.handle)
}

func (h *Handle) Close() error {
	defer func() { h.handle = nil }()

	if h.handle == nil {
		return nil
	}

	return native.Ccodes_handle_delete(h.handle)
}

func handleFinalizer(h *Handle) {
	if h.handle != nil {
		debug.MemoryLeakLogger.Print("handle is not closed")
		h.Close()
	}
}
