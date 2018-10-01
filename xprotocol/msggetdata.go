// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"fmt"

	"github.com/tokublock/tokucore/xbase"
)

// MsgGetData --
type MsgGetData struct {
	InvList []*InvVect
}

// NewMsgGetData -- creates new MsgGetData.
func NewMsgGetData() *MsgGetData {
	return &MsgGetData{}
}

// AddInvVect -- adds an inventory vector to the message.
func (m *MsgGetData) AddInvVect(iv ...*InvVect) error {
	if len(m.InvList)+len(iv) > MaxInvPerMsg {
		return fmt.Errorf("too.many.invvect.in.message[max:%v]", MaxInvPerMsg)
	}
	m.InvList = append(m.InvList, iv...)
	return nil
}

// Encode -- encoding MsgGetData.
func (m *MsgGetData) Encode() []byte {
	buffer := xbase.NewBuffer()
	buffer.WriteVarInt(uint64(len(m.InvList)))
	for _, inv := range m.InvList {
		buffer.WriteU32(uint32(inv.Type))
		buffer.WriteBytes(inv.Hash)
	}
	return buffer.Bytes()
}

// Decode -- decoding data into MsgGetData.
func (m *MsgGetData) Decode(data []byte) error {
	buffer := xbase.NewBufferReader(data)
	count, err := buffer.ReadVarInt()
	if err != nil {
		return err
	}
	if count > MaxInvPerMsg {
		return fmt.Errorf("too.many.invvect.in.message[max:%v]", MaxInvPerMsg)
	}
	for i := 0; i < int(count); i++ {
		var typ uint32
		var hash []byte

		if typ, err = buffer.ReadU32(); err != nil {
			return err
		}
		if hash, err = buffer.ReadBytes(32); err != nil {
			return err
		}
		m.InvList = append(m.InvList, NewInvVect(InvType(typ), hash))
	}
	return nil
}

// Size -- the size of the message.
func (m *MsgGetData) Size() int {
	size := xbase.VarIntSerializeSize(uint64(len(m.InvList)))
	for _, inv := range m.InvList {
		size += inv.Size()
	}
	return size
}

// Command -- returns the protocal command string of this message.
func (m *MsgGetData) Command() string {
	return CommandGetData
}
