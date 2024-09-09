package abin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
)

type nested struct {
	NestedField uint32
}

type test struct {
	ByteField    byte
	ShortField   uint16
	IntField     uint32
	LongField    uint64
	StringField  string
	EmptyString  string
	FloatField   float64
	Nested       nested
	ZeroIntField int32
}

func TestWriteRead(t *testing.T) {
	original := test{
		ByteField:    0x1,
		ShortField:   0x2,
		IntField:     0x3,
		LongField:    0x4,
		StringField:  "test string",
		EmptyString:  "",
		FloatField:   3.14159,
		Nested:       nested{NestedField: 0x5},
		ZeroIntField: 0,
	}

	var buf bytes.Buffer

	err := Write(&buf, binary.LittleEndian, original)
	if err != nil {
		t.Fatalf("Failed to write struct: %v", err)
	}

	var result test

	err = Read(&buf, binary.LittleEndian, &result)
	if err != nil {
		t.Fatalf("Failed to read struct: %v", err)
	}

	if !reflect.DeepEqual(original, result) {
		t.Fatalf("Original and result structs are not equal.\nOriginal: %+v\nResult: %+v", original, result)
	}

	fmt.Printf("%v == %v\n", original, result)
}
