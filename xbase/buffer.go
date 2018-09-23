// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xbase

import (
	"fmt"
)

// Buffer -- a buffer writer/reader for bitcoin.
type Buffer struct {
	seek   int
	length int
	data   []byte
}

// NewBuffer -- creates a new buffer.
func NewBuffer() *Buffer {
	return &Buffer{}
}

// NewBufferReader -- creates a new buffer with data initialize.
func NewBufferReader(data []byte) *Buffer {
	return &Buffer{
		data:   data,
		length: len(data),
	}
}

func (b *Buffer) check(l int) error {
	if (b.seek + l) > b.length {
		return fmt.Errorf("error.seek[%v].required[%v].length[%v]", b.seek, l, b.length)
	}
	return nil
}

// WriteU8 -- write uint8(byte) to buffer.
func (b *Buffer) WriteU8(v uint8) {
	b.data = append(b.data, v)
	b.length++
}

// ReadU8 -- read a byte from the buffer.
func (b *Buffer) ReadU8() (uint8, error) {
	if err := b.check(1); err != nil {
		return 0, err
	}
	v := b.data[b.seek]
	b.seek++
	return v, nil
}

// WriteU16 -- write uint16 with little-endian to the buffer.
func (b *Buffer) WriteU16(v uint32) {
	b.data = append(b.data, byte(v))
	b.data = append(b.data, byte(v>>8))
	b.length += 2
}

// ReadU16 -- read uint16 from the buffer which with the little-endian byteorder.
func (b *Buffer) ReadU16() (uint32, error) {
	if err := b.check(2); err != nil {
		return 0, err
	}

	v := uint32(b.data[b.seek]) |
		uint32(b.data[b.seek+1])<<8
	b.seek += 2
	return v, nil
}

// WriteU32 -- write uint32 with little-endian to the buffer.
func (b *Buffer) WriteU32(v uint32) {
	b.data = append(b.data, byte(v))
	b.data = append(b.data, byte(v>>8))
	b.data = append(b.data, byte(v>>16))
	b.data = append(b.data, byte(v>>24))
	b.length += 4
}

// ReadU32 -- read uint32 from the buffer which with the little-endian byteorder.
func (b *Buffer) ReadU32() (uint32, error) {
	if err := b.check(4); err != nil {
		return 0, err
	}

	v := uint32(b.data[b.seek]) |
		uint32(b.data[b.seek+1])<<8 |
		uint32(b.data[b.seek+2])<<16 |
		uint32(b.data[b.seek+3])<<24
	b.seek += 4
	return v, nil
}

// WriteU64 -- write uint64 to the buffer with little-endian byteorder.
func (b *Buffer) WriteU64(v uint64) {
	b.data = append(b.data, byte(v))
	b.data = append(b.data, byte(v>>8))
	b.data = append(b.data, byte(v>>16))
	b.data = append(b.data, byte(v>>24))
	b.data = append(b.data, byte(v>>32))
	b.data = append(b.data, byte(v>>40))
	b.data = append(b.data, byte(v>>48))
	b.data = append(b.data, byte(v>>56))
	b.length += 8
}

// ReadU64 -- read uint64 from the buffer which with the little-endian byteorder.
func (b *Buffer) ReadU64() (uint64, error) {
	if err := b.check(8); err != nil {
		return 0, err
	}

	v := uint64(b.data[b.seek]) |
		uint64(b.data[b.seek+1])<<8 |
		uint64(b.data[b.seek+2])<<16 |
		uint64(b.data[b.seek+3])<<24 |
		uint64(b.data[b.seek+4])<<32 |
		uint64(b.data[b.seek+5])<<40 |
		uint64(b.data[b.seek+6])<<48 |
		uint64(b.data[b.seek+7])<<56
	b.seek += 8
	return v, nil
}

// WriteVarInt -- write variable int to the buffer.
func (b *Buffer) WriteVarInt(v uint64) {
	switch {
	case v < uint64(0xfd):
		b.WriteU8(uint8(v))
	case v <= 0xffff:
		b.WriteU8(0xfd)
		b.WriteU16(uint32(v))
	case v <= 0xffffffff:
		b.WriteU8(0xfe)
		b.WriteU32(uint32(v))
	default:
		b.WriteU8(0xff)
		b.WriteU64(v)
	}
}

// ReadVarInt -- read variable int from the buffer.
func (b *Buffer) ReadVarInt() (uint64, error) {
	var v uint64
	first, err := b.ReadU8()
	if err != nil {
		return 0, err
	}

	switch first {
	case 0xfd:
		v1, err := b.ReadU16()
		if err != nil {
			return 0, err
		}
		v = uint64(v1)
	case 0xfe:
		v1, err := b.ReadU32()
		if err != nil {
			return 0, err
		}
		v = uint64(v1)
	case 0xff:
		v1, err := b.ReadU64()
		if err != nil {
			return 0, err
		}
		v = v1
	default:
		v = uint64(first)
	}
	return v, nil
}

// VarIntSerializeSize --
// returns the number of bytes it would take to serialize val as a variable length integer.
func VarIntSerializeSize(val uint64) int {
	// The value is small enough to be represented by itself, so it's
	// just 1 byte.
	if val < 0xfd {
		return 1
	}

	// Discriminant 1 byte plus 2 bytes for the uint16.
	if val <= 0xffff {
		return 3
	}

	// Discriminant 1 byte plus 4 bytes for the uint32.
	if val <= 0xffffffff {
		return 5
	}

	// Discriminant 1 byte plus 8 bytes for the uint64.
	return 9
}

// WriteBytes -- write bytes to the buffer.
func (b *Buffer) WriteBytes(v []byte) {
	b.data = append(b.data, v...)
	b.length += len(v)
}

// ReadBytes -- read l-bytes from the buffer.
func (b *Buffer) ReadBytes(l int) ([]byte, error) {
	if l == 0 {
		return nil, nil
	}

	if err := b.check(l); err != nil {
		return nil, err
	}
	v := make([]byte, l)
	copy(v, b.data[b.seek:b.seek+l])
	b.seek += l
	return v, nil
}

// WriteVarBytes -- write variables length bytes to the buffer.
func (b *Buffer) WriteVarBytes(v []byte) {
	len := len(v)
	b.WriteVarInt(uint64(len))
	b.WriteBytes(v)
}

// ReadVarBytes -- read variable length bytes from the buffer.
func (b *Buffer) ReadVarBytes() ([]byte, error) {
	len, err := b.ReadVarInt()
	if err != nil {
		return nil, err
	}
	return b.ReadBytes(int(len))
}

// WriteVarString -- write variables length string to the buffer.
func (b *Buffer) WriteVarString(v string) {
	b.WriteVarBytes([]byte(v))
}

// ReadVarString -- read variable length string from the buffer.
func (b *Buffer) ReadVarString() (string, error) {
	data, err := b.ReadVarBytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Bytes -- returns all the datas in the buffer.
func (b *Buffer) Bytes() []byte {
	return b.data
}

// Remaining -- returns the remaining datas in the buffer.
func (b *Buffer) Remaining() []byte {
	size := b.length - b.seek
	remaining := make([]byte, size)
	copy(remaining[:], b.data[b.seek:b.length])
	return remaining
}

// Seek -- returns the seek postion.
func (b *Buffer) Seek() int {
	return b.seek
}

// Len -- returns the buffer length.
func (b *Buffer) Len() int {
	return b.length
}

// End -- returns whether the seek is end or not.
func (b *Buffer) End() bool {
	return b.seek >= b.length
}

// Reset -- reset the seek, length and the data.
func (b *Buffer) Reset() {
	b.seek = 0
	b.length = 0
	b.data = b.data[0:0]
}
