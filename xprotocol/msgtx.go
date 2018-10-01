// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

// MsgTx --
type MsgTx struct {
	Data []byte
}

// NewMsgTx -- creates new MsgTx.
func NewMsgTx(data []byte) *MsgTx {
	return &MsgTx{Data: data}
}

// Encode -- encoding the message to bitcoin protocol format.
func (m *MsgTx) Encode() []byte {
	return m.Data
}

// Decode -- decoding bytes to the message.
func (m *MsgTx) Decode(data []byte) error {
	m.Data = data
	return nil
}

// Size -- the size of the message.
func (m *MsgTx) Size() int {
	return len(m.Data)
}

// Command -- returns the protocol command of this message.
func (m *MsgTx) Command() string {
	return CommandTx
}
