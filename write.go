package abin

import (
	"encoding/binary"
	"io"
	"reflect"
)

func Write(writer io.Writer, order binary.ByteOrder, data any) (err error) {
	val := reflect.ValueOf(data)

	switch val.Kind() {
	// Encode structs inside structs
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)

			err = Write(writer, order, field.Interface())
			if err != nil {
				return
			}
		}

		return

	// Handle non-fixed length types
	case reflect.String:
		return writeString(writer, order, val.String())

	// Default to normal binary encoding
	default:
		return binary.Write(writer, order, data)
	}
}

func writeString(writer io.Writer, order binary.ByteOrder, str string) (err error) {
	length := uint16(len(str))

	err = binary.Write(writer, order, length)
	if err != nil {
		return
	}

	return binary.Write(writer, order, []byte(str))
}
