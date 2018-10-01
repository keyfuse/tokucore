// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgVerAck(t *testing.T) {
	want := NewMsgVerAck()
	encode := want.Encode()

	got := NewMsgVerAck()
	err := got.Decode(encode)
	assert.Nil(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, CommandVersionAck, got.Command())
	assert.Equal(t, want.Size(), got.Size())
}
