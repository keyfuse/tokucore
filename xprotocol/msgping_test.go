// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xprotocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgPing(t *testing.T) {
	want := NewMsgPing(1234)
	encode := want.Encode()

	got := &MsgPing{}
	err := got.Decode(encode)
	assert.Nil(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, CommandPing, got.Command())
	assert.Equal(t, want.Size(), got.Size())
}
