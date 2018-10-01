// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

// MsgVerAck --
type MsgVerAck struct{}

// NewMsgVerAck -- creates new MsgVerAck.
func NewMsgVerAck() *MsgVerAck {
	return &MsgVerAck{}
}

// Encode -- encoding MsgVersion message to bitcoin protocol format.
func (m *MsgVerAck) Encode() []byte {
	return nil
}

// Decode -- decoding bytes to MsgVerAck message.
func (m *MsgVerAck) Decode(data []byte) error {
	return nil
}

// Size -- the size of the message.
func (m *MsgVerAck) Size() int {
	return len(CommandVersionAck)
}

// Command -- returns the protocol command of this message.
func (m *MsgVerAck) Command() string {
	return CommandVersionAck
}
