// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xprotocol

import (
	"bytes"
	"testing"

	"github.com/keyfuse/tokucore/network"
	"github.com/stretchr/testify/assert"
)

func TestMsgGetHeaders(t *testing.T) {
	network := network.MainNet
	want := NewMsgGetHeaders(network)
	err := want.AddBlockLocatorHash(bytes.Repeat([]byte{0x00}, 32))
	assert.Nil(t, err)
	encode := want.Encode()

	got := NewMsgGetHeaders(network)
	err = got.Decode(encode)
	assert.Nil(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, CommandGetHeaders, got.Command())
	assert.Equal(t, want.Size(), got.Size())
}
