// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgPong(t *testing.T) {
	want := NewMsgPong(1234)
	encode := want.Encode()

	got := &MsgPong{}
	err := got.Decode(encode)
	assert.Nil(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, CommandPong, got.Command())
	assert.Equal(t, want.Size(), got.Size())
}
