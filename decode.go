package binary

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

type decoder struct {
	order  binary.ByteOrder
	reader io.Reader
}

func (d *decoder) decode(v any) (err error) {
	switch v := v.(type) {
	case *bool:
		*v, err = d.bool()
		return

	case *uint8:
		*v, err = d.uint8()
		return
	case *uint16:
		*v, err = d.uint16()
		return
	case *uint32:
		*v, err = d.uint32()
		return
	case *uint64:
		*v, err = d.uint64()
		return
	case *uint:
		*v, err = d.uint()
		return

	case *int8:
		*v, err = d.int8()
		return
	case *int16:
		*v, err = d.int16()
		return
	case *int32:
		*v, err = d.int32()
		return
	case *int64:
		*v, err = d.int64()
		return
	case *int:
		*v, err = d.int()
		return

	case *float32:
		*v, err = d.float32()
		return
	case *float64:
		*v, err = d.float64()
		return

	case *string:
		*v, err = d.string()
		return
	}

	// Otherwise use reflect for structs and slices.
	val := reflect.ValueOf(v).Elem()

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			err = d.decode(val.Field(i).Addr().Interface())
			if err != nil {
				return
			}
		}
		return

	case reflect.Slice:
		var l uint32

		l, err = d.uint32()
		if err != nil {
			return
		}

		val.Set(reflect.MakeSlice(val.Type(), int(l), int(l)))

		for i := 0; i < int(l); i++ {
			err = d.decode(val.Index(i).Addr().Interface())
			if err != nil {
				return
			}
		}
		return
	}

	return fmt.Errorf("unsupported type: %T", v)
}

// Booleans

func (d *decoder) bool() (b bool, err error) {
	var i uint8

	i, err = d.uint8()

	return i != 0, err
}

// Integers

func (d *decoder) uint8() (i uint8, err error) {
	var b [1]byte

	_, err = d.reader.Read(b[:])
	if err != nil {
		return
	}

	return b[0], nil
}

func (d *decoder) uint16() (i uint16, err error) {
	var b [2]byte

	_, err = d.reader.Read(b[:])
	if err != nil {
		return
	}

	return d.order.Uint16(b[:]), nil
}

func (d *decoder) uint32() (i uint32, err error) {
	var b [4]byte

	_, err = d.reader.Read(b[:])
	if err != nil {
		return
	}

	return d.order.Uint32(b[:]), nil
}

func (d *decoder) uint64() (i uint64, err error) {
	var b [8]byte

	_, err = d.reader.Read(b[:])
	if err != nil {
		return
	}

	return d.order.Uint64(b[:]), nil
}

func (d *decoder) uint() (i uint, err error) {
	in, err := d.uint64()

	return uint(in), err
}

func (d *decoder) int8() (i int8, err error) {
	in, err := d.uint8()

	return int8(in), err
}

func (d *decoder) int16() (i int16, err error) {
	in, err := d.uint16()

	return int16(in), err
}

func (d *decoder) int32() (i int32, err error) {
	in, err := d.uint32()

	return int32(in), err
}

func (d *decoder) int64() (i int64, err error) {
	in, err := d.uint64()

	return int64(in), err
}

func (d *decoder) int() (i int, err error) {
	in, err := d.int64()

	return int(in), err
}

// Floats

func (d *decoder) float32() (f float32, err error) {
	in, err := d.uint32()

	return math.Float32frombits(in), err
}

func (d *decoder) float64() (f float64, err error) {
	in, err := d.uint64()

	return math.Float64frombits(in), err
}

// Strings

func (d *decoder) string() (s string, err error) {
	l, err := d.uint32()
	if err != nil {
		return
	}

	b := make([]byte, l)
	_, err = d.reader.Read(b)

	return string(b), err
}
