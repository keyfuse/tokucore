// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"github.com/tokublock/tokucore/xbase"
)

// MsgPong --
type MsgPong struct {
	Nonce uint64
}

// NewMsgPong -- creates new MsgPong.
func NewMsgPong(nonce uint64) *MsgPong {
	return &MsgPong{
		Nonce: nonce,
	}
}

// Encode -- encoding to bitcoin protocol format.
func (m *MsgPong) Encode() []byte {
	buffer := xbase.NewBuffer()
	buffer.WriteU64(m.Nonce)
	return buffer.Bytes()
}

// Decode -- decoding from the bitcoin protocol format.
func (m *MsgPong) Decode(data []byte) error {
	var err error
	buffer := xbase.NewBufferReader(data)

	if m.Nonce, err = buffer.ReadU64(); err != nil {
		return err
	}
	return nil
}

// Size -- the size of the message.
func (m *MsgPong) Size() int {
	return 8
}

// Command -- returns the protocal command string of this message.
func (m *MsgPong) Command() string {
	return CommandPong
}
