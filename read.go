package abin

import (
	"encoding/binary"
	"io"
	"reflect"
)

func Read(reader io.Reader, order binary.ByteOrder, data any) (err error) {
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	switch val.Kind() {
	// Read structs inside structs
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)

			err = Read(reader, order, field.Addr().Interface())
			if err != nil {
				return
			}
		}

		return

	// Read non-fixed length types
	case reflect.String:
		return readString(reader, order, val)

	// Default to normal binary encoding
	default:
		return binary.Read(reader, order, val.Addr().Interface())
	}
}

func readString(reader io.Reader, order binary.ByteOrder, val reflect.Value) (err error) {
	var length uint16

	err = binary.Read(reader, order, &length)
	if err != nil {
		return
	}

	str := make([]byte, length)

	_, err = io.ReadFull(reader, str)
	if err != nil {
		return
	}

	val.SetString(string(str))

	return
}
