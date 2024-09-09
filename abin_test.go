package abin

import (
	"bytes"
	"encoding/binary"
	"testing"
)

type example struct {
	A int64
	B string
	C bool
	D uint64
	E float64
	F string
	G string
	H uint16
}

func TestWriteRead(t *testing.T) {
	in := example{
		A: 123,
		B: "hello",
		C: true,
		D: 1234567890,
		E: 123.456,
		F: "joe",
		G: "mama",
		H: 12345,
	}

	var buf bytes.Buffer

	err := Write(&buf, binary.LittleEndian, in)
	if err != nil {
		t.Fatalf("error writing: %v", err)
	}

	var out example

	err = Read(&buf, binary.LittleEndian, &out)
	if err != nil {
		t.Fatalf("error reading: %v", err)
	}

	if in.A != out.A {
		t.Errorf("A: want %d, got %d", in.A, out.A)
	}

	if in.B != out.B {
		t.Errorf("B: want '%s', got '%s'", in.B, out.B)
	}

	if in.C != out.C {
		t.Errorf("C: want %v, got %v", in.C, out.C)
	}

	if in.D != out.D {
		t.Errorf("D: want %d, got %d", in.D, out.D)
	}

	if in.E != out.E {
		t.Errorf("E: want %f, got %f", in.E, out.E)
	}

	if in.F != out.F {
		t.Errorf("F: want '%s', got '%s'", in.F, out.F)
	}

	if in.G != out.G {
		t.Errorf("G: want '%s', got '%s'", in.G, out.G)
	}

	if in.H != out.H {
		t.Errorf("H: want %d, got %d", in.H, out.H)
	}
}
