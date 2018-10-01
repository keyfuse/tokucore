// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

// MsgUnhandle --
type MsgUnhandle struct {
	command string
}

// NewMsgUnhandle -- creates new MsgUnhandle.
func NewMsgUnhandle(command string) *MsgUnhandle {
	return &MsgUnhandle{
		command: command,
	}
}

// Encode -- encoding the message to bitcoin protocol format.
func (m *MsgUnhandle) Encode() []byte {
	return nil
}

// Decode -- decoding bytes to the message.
func (m *MsgUnhandle) Decode(data []byte) error {
	return nil
}

// Size -- the size of the message.
func (m *MsgUnhandle) Size() int {
	return len(m.command)
}

// Command -- returns the protocol command of this message.
func (m *MsgUnhandle) Command() string {
	return m.command
}
