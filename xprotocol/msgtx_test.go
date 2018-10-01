// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgTx(t *testing.T) {
	want := NewMsgTx([]byte{0x01})
	encode := want.Encode()

	got := &MsgTx{}
	err := got.Decode(encode)
	assert.Nil(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, CommandTx, got.Command())
	assert.Equal(t, want.Size(), got.Size())
}
