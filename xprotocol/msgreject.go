// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"github.com/tokublock/tokucore/xbase"
)

// RejectCode --
type RejectCode uint8

//
const (
	RejectMalformed       RejectCode = 0x01
	RejectInvalid         RejectCode = 0x10
	RejectObsolete        RejectCode = 0x11
	RejectDuplicate       RejectCode = 0x12
	RejectNonstandard     RejectCode = 0x40
	RejectDust            RejectCode = 0x41
	RejectInsufficientFee RejectCode = 0x42
	RejectCheckpoint      RejectCode = 0x43
)

// MsgReject --
type MsgReject struct {
	Cmd    string
	Code   uint8
	Reason string
	Hash   []byte
}

// NewMsgReject -- creates new MsgReject.
func NewMsgReject(command string, code RejectCode, reason string) *MsgReject {
	return &MsgReject{
		Cmd:    command,
		Code:   uint8(code),
		Reason: reason,
	}
}

// Encode -- encoding to bitcoin protocol format.
func (m *MsgReject) Encode() []byte {
	buffer := xbase.NewBuffer()
	buffer.WriteVarString(m.Cmd)
	buffer.WriteU8(m.Code)
	buffer.WriteVarString(m.Reason)
	if m.Cmd == CommandBlock || m.Cmd == CommandTx {
		buffer.WriteBytes(m.Hash)
	}
	return buffer.Bytes()
}

// Decode -- decoding from the bitcoin protocol format.
func (m *MsgReject) Decode(data []byte) error {
	var err error
	buffer := xbase.NewBufferReader(data)

	if m.Cmd, err = buffer.ReadVarString(); err != nil {
		return err
	}
	if m.Code, err = buffer.ReadU8(); err != nil {
		return err
	}
	if m.Reason, err = buffer.ReadVarString(); err != nil {
		return err
	}
	if m.Cmd == CommandBlock || m.Cmd == CommandTx {
		if m.Hash, err = buffer.ReadBytes(32); err != nil {
			return err
		}
	}
	return nil
}

// Size -- the size of the message.
func (m *MsgReject) Size() int {
	return len(m.Cmd) + 1 + len(m.Reason) + len(m.Hash)
}

// Command -- returns the protocal command string of this message.
func (m *MsgReject) Command() string {
	return CommandReject
}
