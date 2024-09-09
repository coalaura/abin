package binary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

type test struct {
	Num   int64
	Str   string
	Float float64
	Bool  bool
}

func TestBinary(t *testing.T) {
	original := test{
		Num:   42,
		Str:   "hello world",
		Float: 3.1459,
		Bool:  true,
	}

	var buf bytes.Buffer

	t.Log("writing")
	if err := Write(&buf, binary.BigEndian, original); err != nil {
		t.Fatal(err)
	}

	var result test

	t.Log("reading")
	reader := bytes.NewReader(buf.Bytes())

	if err := Read(reader, binary.BigEndian, &result); err != nil {
		t.Fatal(err)
	}

	if result.Num != original.Num {
		t.Fatalf("expected %d, got %d", original.Num, result.Num)
	}

	if result.Str != original.Str {
		t.Fatalf("expected '%s', got '%s'", original.Str, result.Str)
	}

	if result.Float != original.Float {
		t.Fatalf("expected %f, got %f", original.Float, result.Float)
	}

	if result.Bool != original.Bool {
		t.Fatalf("expected %t, got %t", original.Bool, result.Bool)
	}

	fmt.Printf("%v == %v\n", original, result)
}
