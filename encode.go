package binary

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

type encoder struct {
	order  binary.ByteOrder
	writer io.Writer
}

func (e *encoder) encode(v any) (err error) {
	switch v := v.(type) {
	case bool:
		return e.bool(v)

	case uint8:
		return e.uint8(v)
	case uint16:
		return e.uint16(v)
	case uint32:
		return e.uint32(v)
	case uint64:
		return e.uint64(v)
	case uint:
		return e.uint(v)

	case int8:
		return e.int8(v)
	case int16:
		return e.int16(v)
	case int32:
		return e.int32(v)
	case int64:
		return e.int64(v)
	case int:
		return e.int(v)

	case float32:
		return e.float32(v)
	case float64:
		return e.float64(v)

	case string:
		return e.string(v)
	}

	// Otherwise use reflect for structs and slices.
	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if err = e.encode(val.Field(i).Interface()); err != nil {
				return
			}
		}
		return

	case reflect.Slice:
		err = e.encode(uint32(val.Len()))
		if err != nil {
			return
		}

		for i := 0; i < val.Len(); i++ {
			if err = e.encode(val.Index(i).Interface()); err != nil {
				return
			}
		}
		return
	}

	return fmt.Errorf("unsupported type: %T", v)
}

// Booleans

func (e *encoder) bool(b bool) (err error) {
	if b {
		_, err = e.writer.Write([]byte{1})
	} else {
		_, err = e.writer.Write([]byte{0})
	}
	return
}

// Integers

func (e *encoder) uint8(i uint8) (err error) {
	_, err = e.writer.Write([]byte{i})
	return
}

func (e *encoder) uint16(i uint16) (err error) {
	if e.order == binary.BigEndian {
		_, err = e.writer.Write([]byte{
			byte(i >> 8),
			byte(i),
		})
	} else {
		_, err = e.writer.Write([]byte{
			byte(i),
			byte(i >> 8),
		})
	}
	return
}

func (e *encoder) uint32(i uint32) (err error) {
	if e.order == binary.BigEndian {
		_, err = e.writer.Write([]byte{
			byte(i >> 24),
			byte(i >> 16),
			byte(i >> 8),
			byte(i),
		})
	} else {
		_, err = e.writer.Write([]byte{
			byte(i),
			byte(i >> 8),
			byte(i >> 16),
			byte(i >> 24),
		})
	}
	return
}

func (e *encoder) uint64(i uint64) (err error) {
	if e.order == binary.BigEndian {
		_, err = e.writer.Write([]byte{
			byte(i >> 56),
			byte(i >> 48),
			byte(i >> 40),
			byte(i >> 32),
			byte(i >> 24),
			byte(i >> 16),
			byte(i >> 8),
			byte(i),
		})
	} else {
		_, err = e.writer.Write([]byte{
			byte(i),
			byte(i >> 8),
			byte(i >> 16),
			byte(i >> 24),
			byte(i >> 32),
			byte(i >> 40),
			byte(i >> 48),
			byte(i >> 56),
		})
	}
	return
}

func (e *encoder) uint(i uint) (err error) {
	return e.uint64(uint64(i))
}

func (e *encoder) int8(i int8) (err error) {
	return e.uint8(uint8(i))
}

func (e *encoder) int16(i int16) (err error) {
	return e.uint16(uint16(i))
}

func (e *encoder) int32(i int32) (err error) {
	return e.uint32(uint32(i))
}

func (e *encoder) int64(i int64) (err error) {
	return e.uint64(uint64(i))
}

func (e *encoder) int(i int) (err error) {
	return e.int64(int64(i))
}

// Floats

func (e *encoder) float32(f float32) (err error) {
	return e.uint32(math.Float32bits(f))
}

func (e *encoder) float64(f float64) (err error) {
	return e.uint64(math.Float64bits(f))
}

// Strings

func (e *encoder) string(s string) (err error) {
	l := uint32(len(s))
	if err = e.uint32(l); err != nil {
		return
	}

	_, err = e.writer.Write([]byte(s))
	return
}
