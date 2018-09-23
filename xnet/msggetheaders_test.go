// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xnet

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgGetHeaders(t *testing.T) {
	want := NewMsgGetHeaders(MainNet)
	want.AddBlockLocatorHash(bytes.Repeat([]byte{0x00}, 32))
	encode := want.Encode()

	got := NewMsgGetHeaders(MainNet)
	err := got.Decode(encode)
	assert.Nil(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, CommandGetHeaders, got.Command())
	assert.Equal(t, want.Size(), got.Size())
}
