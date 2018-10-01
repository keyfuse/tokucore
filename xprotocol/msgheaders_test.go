// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
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
