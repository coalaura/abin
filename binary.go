package binary

import (
	"encoding/binary"
	"io"
)

var (
	BigEndian    = binary.BigEndian
	LittleEndian = binary.LittleEndian
)

func Write(w io.Writer, order binary.ByteOrder, v any) (err error) {
	e := encoder{
		order:  order,
		writer: w,
	}

	return e.encode(v)
}

func Read(r io.Reader, order binary.ByteOrder, v any) (err error) {
	d := decoder{
		order:  order,
		reader: r,
	}

	return d.decode(v)
}
