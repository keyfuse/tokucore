// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"github.com/tokublock/tokucore/xbase"
)

// MsgPing --
type MsgPing struct {
	Nonce uint64
}

// NewMsgPing -- creates new MsgPing.
func NewMsgPing(nonce uint64) *MsgPing {
	return &MsgPing{
		Nonce: nonce,
	}
}

// Encode -- encoding to bitcoin protocol format.
func (m *MsgPing) Encode() []byte {
	buffer := xbase.NewBuffer()
	buffer.WriteU64(m.Nonce)
	return buffer.Bytes()
}

// Decode -- decoding from the bitcoin protocol format.
func (m *MsgPing) Decode(data []byte) error {
	var err error
	buffer := xbase.NewBufferReader(data)

	if m.Nonce, err = buffer.ReadU64(); err != nil {
		return err
	}
	return nil
}

// Size -- the size of the message.
func (m *MsgPing) Size() int {
	return 8
}

// Command -- returns the protocal command string of this message.
func (m *MsgPing) Command() string {
	return CommandPing
}
