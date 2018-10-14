// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokublock/tokucore/xbase"
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

func TestMsgRejectError(t *testing.T) {
	m := NewMsgReject("tx", RejectMalformed, "xx")

	f0 := func(buffer *xbase.Buffer) {
		buffer.WriteVarString((m.Cmd))
	}

	f1 := func(buffer *xbase.Buffer) {
		buffer.WriteU8((m.Code))
	}

	f2 := func(buffer *xbase.Buffer) {
		buffer.WriteVarString((m.Reason))
	}

	buffer := xbase.NewBuffer()
	fs := []func(buff *xbase.Buffer){f0, f1, f2}
	for _, fn := range fs {
		msg := &MsgReject{}
		err := msg.Decode(buffer.Bytes())
		assert.NotNil(t, err)
		fn(buffer)
	}
}
