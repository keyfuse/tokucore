// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"fmt"

	"github.com/tokublock/tokucore/xbase"
)

// MsgInv --
type MsgInv struct {
	Invs []*InvVect
}

// NewMsgInv -- creates new MsgInv.
func NewMsgInv() *MsgInv {
	return &MsgInv{}
}

// AddInvVect -- adds an inventory vector to the message.
func (m *MsgInv) AddInvVect(iv ...*InvVect) error {
	if len(m.Invs)+1 > MaxInvPerMsg {
		return fmt.Errorf("too.many.invvect.in.message[max:%v]", MaxInvPerMsg)
	}
	m.Invs = append(m.Invs, iv...)
	return nil
}

// Encode -- encoding MsgInv.
func (m *MsgInv) Encode() []byte {
	buffer := xbase.NewBuffer()
	buffer.WriteVarInt(uint64(len(m.Invs)))
	for _, inv := range m.Invs {
		buffer.WriteU32(uint32(inv.Type))
		buffer.WriteBytes(inv.Hash)
	}
	return buffer.Bytes()
}

// Decode -- decoding data to MsgInv.
func (m *MsgInv) Decode(data []byte) error {
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
		m.Invs = append(m.Invs, NewInvVect(InvType(typ), hash))
	}
	return nil
}

// Size -- the size of the message.
func (m *MsgInv) Size() int {
	size := xbase.VarIntSerializeSize(uint64(len(m.Invs)))
	for _, inv := range m.Invs {
		size += inv.Size()
	}
	return size
}

// Command -- returns the protocol command of this message.
func (m *MsgInv) Command() string {
	return CommandInventory
}
