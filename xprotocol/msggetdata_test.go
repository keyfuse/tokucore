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

func TestMsgGetData(t *testing.T) {
	want := NewMsgGetData()
	want.AddInvVect(&InvVect{Type: InvTypeTx, Hash: bytes.Repeat([]byte{0x00}, 32)})
	encode := want.Encode()

	got := NewMsgGetData()
	err := got.Decode(encode)
	assert.Nil(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, CommandGetData, got.Command())
	assert.Equal(t, want.Size(), got.Size())
}
