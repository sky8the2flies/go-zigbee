package buffer

import (
	"bytes"
	"encoding/binary"
)

type Scannable struct {
	buffer *bytes.Reader
	raw    []byte
}

func NewScannable(data []byte) *Scannable {
	return &Scannable{buffer: bytes.NewReader(data), raw: data}
}

func (dv *Scannable) GetUint16(offset int, littleEndian bool) (uint16, error) {
	dv.buffer.Seek(int64(offset), 0) // Seek to offset
	var value uint16
	var order binary.ByteOrder
	if littleEndian {
		order = binary.LittleEndian
	} else {
		order = binary.BigEndian
	}
	err := binary.Read(dv.buffer, order, &value)
	return value, err
}

func (dv *Scannable) GetUint8(offset int, littleEndian bool) (uint8, error) {
	dv.buffer.Seek(int64(offset), 0) // Seek to offset
	var value uint8
	var order binary.ByteOrder
	if littleEndian {
		order = binary.LittleEndian
	} else {
		order = binary.BigEndian
	}
	err := binary.Read(dv.buffer, order, &value)
	return value, err
}

func (dv *Scannable) Raw() []byte {
	return dv.raw
}
