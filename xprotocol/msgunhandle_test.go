// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgUnhandle(t *testing.T) {
	want := NewMsgUnhandle("xx")
	encode := want.Encode()

	got := NewMsgUnhandle("xx")
	err := got.Decode(encode)
	assert.Nil(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, "xx", got.Command())
	assert.Equal(t, want.Size(), got.Size())
}
