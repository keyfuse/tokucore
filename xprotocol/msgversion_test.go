// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xbase"
)

func TestMsgVersion(t *testing.T) {
	network := network.MainNet
	want := NewMsgVersion(network)
	encode := want.Encode()

	got := NewMsgVersion(network)
	err := got.Decode(encode)
	assert.Nil(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, CommandVersion, got.Command())
	assert.Equal(t, want.Size(), got.Size())
}

func TestMsgVersionError(t *testing.T) {
	network := network.MainNet
	m := NewMsgVersion(network)

	f0 := func(buffer *xbase.Buffer) {
		buffer.WriteU32((m.Version))
	}

	f1 := func(buffer *xbase.Buffer) {
		buffer.WriteU64((m.Services))
	}

	f2 := func(buffer *xbase.Buffer) {
		buffer.WriteU64(m.Timestamp)
	}

	f3 := func(buffer *xbase.Buffer) {
		buffer.WriteU64(m.ServicesYou)
	}

	f4 := func(buffer *xbase.Buffer) {
		buffer.WriteBytes(m.AddressYou)
	}

	f5 := func(buffer *xbase.Buffer) {
		buffer.WriteU16(m.PortYou)
	}

	f6 := func(buffer *xbase.Buffer) {
		buffer.WriteU64(m.ServicesMe)
	}

	f7 := func(buffer *xbase.Buffer) {
		buffer.WriteBytes(m.AddressMe)
	}

	f8 := func(buffer *xbase.Buffer) {
		buffer.WriteU16(m.PortMe)
	}

	f9 := func(buffer *xbase.Buffer) {
		buffer.WriteU64(m.Nonce)
	}

	f10 := func(buffer *xbase.Buffer) {
		buffer.WriteVarString(m.UserAgent)
	}

	f11 := func(buffer *xbase.Buffer) {
		buffer.WriteU32((m.LastBlock))
	}

	f12 := func(buffer *xbase.Buffer) {
		buffer.WriteU8(m.Relay)
	}

	buffer := xbase.NewBuffer()
	fs := []func(buff *xbase.Buffer){f0, f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11, f12}
	for _, fn := range fs {
		msg := &MsgVersion{}
		err := msg.Decode(buffer.Bytes())
		assert.NotNil(t, err)
		fn(buffer)
	}
}
