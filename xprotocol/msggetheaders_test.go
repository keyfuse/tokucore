// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokublock/tokucore/network"
)

func TestMsgGetHeaders(t *testing.T) {
	network := network.MainNet
	want := NewMsgGetHeaders(network)
	want.AddBlockLocatorHash(bytes.Repeat([]byte{0x00}, 32))
	encode := want.Encode()

	got := NewMsgGetHeaders(network)
	err := got.Decode(encode)
	assert.Nil(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, CommandGetHeaders, got.Command())
	assert.Equal(t, want.Size(), got.Size())
}
