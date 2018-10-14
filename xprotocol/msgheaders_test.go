// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokublock/tokucore/xbase"
)

func TestMsgHeaders(t *testing.T) {
	want := NewMsgHeaders()
	want.AddBlockHeader(&BlockHeader{
		Version:    1,
		PrevBlock:  bytes.Repeat([]byte{0x00}, 32),
		MerkleRoot: bytes.Repeat([]byte{0x00}, 32),
		Timestamp:  999,
		Nonce:      888,
	})
	encode := want.Encode()

	got := NewMsgHeaders()
	err := got.Decode(encode)
	assert.Nil(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, CommandHeaders, got.Command())
	assert.Equal(t, want.Size(), got.Size())
}

func TestMsgHeadersError(t *testing.T) {
	want := NewMsgHeaders()
	for i := 0; i < MaxBlockHeadersPerMsg; i++ {
		want.AddBlockHeader(&BlockHeader{})
	}
	err := want.AddBlockHeader(&BlockHeader{})
	assert.NotNil(t, err)
}

func TestMsgHeadersError1(t *testing.T) {
	hdr := &BlockHeader{}

	f0 := func(buffer *xbase.Buffer) {
		buffer.WriteVarInt(uint64(1))
	}

	f1 := func(buffer *xbase.Buffer) {
		buffer.WriteU32(hdr.Version)
	}

	f2 := func(buffer *xbase.Buffer) {
		buffer.WriteBytes(hdr.PrevBlock)
	}

	f3 := func(buffer *xbase.Buffer) {
		buffer.WriteBytes(hdr.MerkleRoot)
	}

	f4 := func(buffer *xbase.Buffer) {
		buffer.WriteU32(hdr.Timestamp)
	}

	f5 := func(buffer *xbase.Buffer) {
		buffer.WriteU32(hdr.Bits)
	}

	f6 := func(buffer *xbase.Buffer) {
		buffer.WriteU32(hdr.Nonce)
	}

	buffer := xbase.NewBuffer()
	fs := []func(buff *xbase.Buffer){f0, f1, f2, f3, f4, f5, f6}
	for _, fn := range fs {
		msg := NewMsgHeaders()
		err := msg.Decode(buffer.Bytes())
		assert.NotNil(t, err)
		fn(buffer)
	}
}
