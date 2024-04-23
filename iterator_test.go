package codes_test

import (
	"testing"

	codes "github.com/scorix/go-eccodes"
	"github.com/scorix/go-eccodes/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIterator_Next(t *testing.T) {
	files := []string{
		"test-data/dirpw_surface_1.grib2",
		"test-data/prate_surface_0.grib2",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			f, err := io.OpenFile(file, "r")
			require.NoError(t, err)
			defer f.Close()

			grib, err := codes.OpenFile(f)
			require.NoError(t, err)
			defer grib.Close()

			handle, err := grib.Handle()
			require.NoError(t, err)
			defer handle.Close()

			msg := handle.Message()
			defer msg.Close()

			i, err := msg.Iterator()
			require.NoError(t, err)
			defer i.Close()

			var cnt int
			for i.HasNext() {
				i.Next()
				cnt++
			}

			assert.Equal(t, 1440*721, cnt)
		})
	}
}
