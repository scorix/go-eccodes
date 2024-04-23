package codes

import (
	"errors"
	"math"
	"runtime"
	"slices"

	"github.com/scorix/go-eccodes/debug"
	"github.com/scorix/go-eccodes/native"
)

type Message interface {
	isOpen() bool

	GetString(key string) (string, error)

	GetLong(key string) (int64, error)
	SetLong(key string, value int64) error

	GetDouble(key string) (float64, error)
	SetDouble(key string, value float64) error

	Data() (latitudes []float64, longitudes []float64, values []float64, err error)
	DataUnsafe() (latitudes *Float64ArrayUnsafe, longitudes *Float64ArrayUnsafe, values *Float64ArrayUnsafe, err error)

	NearestFind(latitude float64, longitude float64) (outLat float64, outLon float64, value float64, distance float64, index int32, err error)

	Iterator() (Iterator, error)

	Close() error
}

type message struct {
	handle native.Ccodes_handle
}

func newMessage(h native.Ccodes_handle) Message {
	m := &message{handle: h}
	runtime.SetFinalizer(m, messageFinalizer)

	// set missing value to NaN
	m.SetDouble(parameterMissingValue, math.NaN()) //nolint: errcheck

	return m
}

func (m *message) isOpen() bool {
	return m.handle != nil
}

func (m *message) GetString(key string) (string, error) {
	return native.Ccodes_get_string(m.handle, key)
}

func (m *message) GetLong(key string) (int64, error) {
	return native.Ccodes_get_long(m.handle, key)
}

func (m *message) SetLong(key string, value int64) error {
	return native.Ccodes_set_long(m.handle, key, value)
}

func (m *message) GetDouble(key string) (float64, error) {
	return native.Ccodes_get_double(m.handle, key)
}

func (m *message) SetDouble(key string, value float64) error {
	return native.Ccodes_set_double(m.handle, key, value)
}

func (m *message) Data() (latitudes []float64, longitudes []float64, values []float64, err error) {
	return native.Ccodes_grib_get_data(m.handle)
}

func (m *message) NearestFind(latitude float64, longitude float64) (outLat float64, outLon float64, value float64, distance float64, index int32, err error) {
	n, err := native.Ccodes_grib_nearest_new(m.handle)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	defer native.Ccodes_grib_nearest_delete(n) //nolint: errcheck

	lats, lons, values, distances, indexes, err := native.Ccodes_grib_nearest_find(n, m.handle, latitude, longitude)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	idxs := slices.DeleteFunc(indexes, func(a int32) bool {
		return a == 0
	})

	if len(idxs) == 0 {
		return 0, 0, 0, 0, 0, errors.New("no nearest found")
	}

	minDistance := distances[0]
	minIdx := int32(0)
	for i := 1; i < 4; i++ {
		if distances[i] < minDistance {
			minIdx = int32(i)
			minDistance = distances[i]
		}
	}

	return lats[minIdx], lons[minIdx], values[minIdx], minDistance, minIdx, nil
}

func (m *message) DataUnsafe() (latitudes *Float64ArrayUnsafe, longitudes *Float64ArrayUnsafe, values *Float64ArrayUnsafe, err error) {
	lats, lons, vals, err := native.Ccodes_grib_get_data_unsafe(m.handle)
	if err != nil {
		return nil, nil, nil, err
	}

	return newFloat64ArrayUnsafe(lats), newFloat64ArrayUnsafe(lons), newFloat64ArrayUnsafe(vals), nil
}

func (m *message) Iterator() (Iterator, error) {
	if m.handle == nil {
		return nil, errors.New("handle has been closed")
	}

	return newIterator(m.handle)
}

func (m *message) Close() error {
	defer func() { m.handle = nil }()

	return nil
}

func messageFinalizer(m *message) {
	if m.isOpen() {
		debug.MemoryLeakLogger.Print("message is not closed")
		m.Close()
	}
}
