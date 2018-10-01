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

func TestMsgReject(t *testing.T) {
	want := NewMsgReject("tx", RejectMalformed, "xx")
	want.Hash = bytes.Repeat([]byte{0x01}, 32)
	encode := want.Encode()

	got := &MsgReject{}
	err := got.Decode(encode)
	assert.Nil(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, CommandReject, got.Command())
	assert.Equal(t, want.Size(), got.Size())
}
