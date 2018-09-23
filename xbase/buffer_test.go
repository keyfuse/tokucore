// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xbase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	writer := NewBuffer()
	writer.WriteU8(0)
	writer.WriteBytes([]byte{0x01, 0x02})
	writer.WriteU16(65534)
	writer.WriteU32(653400000)
	writer.WriteU64(uint64(43000000000))

	datas := writer.Bytes()
	reader := NewBufferReader(datas)

	// U8.
	{
		v, err := reader.ReadU8()
		assert.Nil(t, err)
		want := uint8(0)
		assert.Equal(t, want, v)
	}

	// bytes.
	{
		v, err := reader.ReadBytes(2)
		assert.Nil(t, err)
		want := []byte{0x01, 0x02}
		assert.Equal(t, want, v)
	}

	// U16.
	{
		end := reader.End()
		assert.False(t, end)

		v, err := reader.ReadU16()
		assert.Nil(t, err)
		want := uint32(65534)
		assert.Equal(t, want, v)
	}

	// Seek, Len.
	{
		seek := reader.Seek()
		assert.Equal(t, 5, seek)

		pos := reader.Len()
		assert.Equal(t, 17, pos)
	}

	// U32.
	{
		end := reader.End()
		assert.False(t, end)

		v, err := reader.ReadU32()
		assert.Nil(t, err)
		want := uint32(653400000)
		assert.Equal(t, want, v)
	}

	// U64.
	{

		v, err := reader.ReadU64()
		assert.Nil(t, err)
		want := uint64(43000000000)
		assert.Equal(t, want, v)
	}

	// Error.
	{
		end := reader.End()
		assert.True(t, end)

		_, err := reader.ReadU8()
		assert.NotNil(t, err)

		_, err = reader.ReadBytes(1)
		assert.NotNil(t, err)

		_, err = reader.ReadU16()
		assert.NotNil(t, err)

		_, err = reader.ReadU32()
		assert.NotNil(t, err)
	}

	// Reset.
	{
		writer.Reset()
	}
}

func TestBufferVarInt(t *testing.T) {
	writer := NewBuffer()
	writer.WriteVarInt(uint64(8))
	writer.WriteVarInt(uint64(0xFD + 1))
	writer.WriteVarInt(uint64(0xFFFF + 1))
	writer.WriteVarInt(uint64(0xFFFFFFFF + 1))

	datas := writer.Bytes()
	reader := NewBufferReader(datas)

	// <0xFD
	{
		v, err := reader.ReadVarInt()
		assert.Nil(t, err)
		want := uint64(8)
		assert.Equal(t, want, v)

		size := VarIntSerializeSize(want)
		assert.Equal(t, 1, size)
	}

	// >0xFD
	{
		v, err := reader.ReadVarInt()
		assert.Nil(t, err)
		want := uint64(0xFD + 1)
		assert.Equal(t, want, v)

		size := VarIntSerializeSize(want)
		assert.Equal(t, 3, size)
	}

	// >uint16
	{
		v, err := reader.ReadVarInt()
		assert.Nil(t, err)
		want := uint64(0xFFFF + 1)
		assert.Equal(t, want, v)

		size := VarIntSerializeSize(want)
		assert.Equal(t, 5, size)
	}

	// >uint32
	{
		v, err := reader.ReadVarInt()
		assert.Nil(t, err)
		want := uint64(0xFFFFFFFF + 1)
		assert.Equal(t, want, v)

		size := VarIntSerializeSize(want)
		assert.Equal(t, 9, size)
	}
}

func TestBufferVarBytes(t *testing.T) {
	writer := NewBuffer()
	writer.WriteVarInt(uint64(8))
	writer.WriteVarBytes([]byte{0x01, 0x02})
	writer.WriteVarBytes([]byte{0x03, 0x04})

	datas := writer.Bytes()
	reader := NewBufferReader(datas)

	// <0xFD
	{
		v, err := reader.ReadVarInt()
		assert.Nil(t, err)
		want := uint64(8)
		assert.Equal(t, want, v)
	}

	// bytes.
	{
		v, err := reader.ReadVarBytes()
		assert.Nil(t, err)
		want := []byte{0x01, 0x02}
		assert.Equal(t, want, v)
	}

	// bytes.
	{
		v, err := reader.ReadVarBytes()
		assert.Nil(t, err)
		want := []byte{0x03, 0x04}
		assert.Equal(t, want, v)
	}
}

func TestBufferVarString(t *testing.T) {
	writer := NewBuffer()
	writer.WriteVarString("xx")
	writer.WriteBytes([]byte{0x01, 0x02})

	datas := writer.Bytes()
	reader := NewBufferReader(datas)

	v, err := reader.ReadVarString()
	assert.Nil(t, err)
	assert.Equal(t, "xx", v)

	remaing := reader.Remaining()
	assert.Equal(t, []byte{0x01, 0x02}, remaing)
}

func BenchmarkBuffer(b *testing.B) {
	writer := NewBuffer()
	for i := 0; i < b.N; i++ {
		writer.WriteVarInt(uint64(8))
		writer.WriteVarBytes([]byte{0x01, 0x02})
	}

	for i := 0; i < b.N; i++ {
		writer.ReadVarInt()
		writer.ReadVarBytes()
	}
}
