package codes_test

import (
	"testing"

	codes "github.com/scorix/go-eccodes"
	"github.com/scorix/go-eccodes/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage_NearestFind(t *testing.T) {
	f, err := io.OpenFile("test-data/dirpw_surface_1.grib2", "r")
	require.NoError(t, err)
	defer f.Close()

	grib, err := codes.OpenFile(f)
	require.NoError(t, err)

	msg, err := grib.Next()
	require.NoError(t, err)
	defer msg.Close()

	lat, lon, val, _, _, err := msg.NearestFind(77.25, 10)
	require.NoError(t, err)

	assert.Equal(t, float32(77.25), float32(lat))
	assert.Equal(t, float32(10), float32(lon))
	assert.Equal(t, float32(206.98), float32(val))
}

func TestMessage_GetString(t *testing.T) {
	t.Run("without index", func(t *testing.T) {
		f, err := io.OpenFile("test-data/dirpw_surface_1.grib2", "r")
		require.NoError(t, err)
		defer f.Close()

		grib, err := codes.OpenFile(f)
		require.NoError(t, err)

		msg, err := grib.Next()
		require.NoError(t, err)
		defer msg.Close()

		s, err := msg.GetString("time")
		require.NoError(t, err)

		assert.Equal(t, "1800", s)
	})

	t.Run("with index", func(t *testing.T) {
		f, err := codes.OpenFileByPathWithFilter("test-data/dirpw_surface_1.grib2", map[string]any{"shortName": "dirpw"})
		require.NoError(t, err)
		defer f.Close()

		msg, err := f.Next()
		require.NoError(t, err)
		defer msg.Close()

		s, err := msg.GetString("time")
		require.NoError(t, err)

		assert.Equal(t, "1800", s)
	})
}

func BenchmarkMessage_GetString(b *testing.B) {
	f, err := io.OpenFile("test-data/dirpw_surface_1.grib2", "r")
	require.NoError(b, err)
	defer f.Close()

	grib, err := codes.OpenFile(f)
	require.NoError(b, err)

	msg, err := grib.Next()
	require.NoError(b, err)
	defer msg.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.GetString("shortName")
	}
}

func BenchmarkNearestFind(b *testing.B) {
	b.Run("without index", func(b *testing.B) {
		f, err := io.OpenFile("test-data/dirpw_surface_1.grib2", "r")
		require.NoError(b, err)
		defer f.Close()

		grib, err := codes.OpenFile(f)
		require.NoError(b, err)

		msg, err := grib.Next()
		require.NoError(b, err)
		defer msg.Close()

		for i := 0; i < b.N; i++ {
			msg.NearestFind(77.25, 10)
		}
	})

	b.Run("with index", func(b *testing.B) {
		grib, err := codes.OpenFileByPathWithFilter("test-data/dirpw_surface_1.grib2", map[string]any{"shortName": "dirpw"})
		require.NoError(b, err)
		defer grib.Close()

		msg, err := grib.Next()
		require.NoError(b, err)
		defer msg.Close()

		for i := 0; i < b.N; i++ {
			msg.NearestFind(77.25, 10)
		}
	})
}

func BenchmarkData(b *testing.B) {
	f, err := io.OpenFile("test-data/dirpw_surface_1.grib2", "r")
	require.NoError(b, err)
	defer f.Close()

	grib, err := codes.OpenFile(f)
	require.NoError(b, err)

	msg, err := grib.Next()
	require.NoError(b, err)
	defer msg.Close()

	for i := 0; i < b.N; i++ {
		msg.Data()
	}
}
