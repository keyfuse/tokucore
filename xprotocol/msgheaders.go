// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"fmt"

	"github.com/tokublock/tokucore/xbase"
)

// MaxBlockHeadersPerMsg  -- the maximum number of block headers that can be in
// a single bitcoin headers message.
const MaxBlockHeadersPerMsg = 2000

// MsgHeaders --
type MsgHeaders struct {
	Headers []*BlockHeader
}

// NewMsgHeaders -- creates new MsgHeaders.
func NewMsgHeaders() *MsgHeaders {
	return &MsgHeaders{}
}

// AddBlockHeader -- adds a new block header to the message.
func (m *MsgHeaders) AddBlockHeader(headers ...*BlockHeader) error {
	if len(m.Headers)+1 > MaxBlockHeadersPerMsg {
		return fmt.Errorf("too.many.block.headers.in.message[max:%v]", MaxBlockHeadersPerMsg)
	}
	m.Headers = append(m.Headers, headers...)
	return nil
}

// Encode -- encoding to bitcoin protocol format.
func (m *MsgHeaders) Encode() []byte {
	buffer := xbase.NewBuffer()
	buffer.WriteVarInt(uint64(len(m.Headers)))
	for _, hdr := range m.Headers {
		buffer.WriteU32(hdr.Version)
		buffer.WriteBytes(hdr.PrevBlock)
		buffer.WriteBytes(hdr.MerkleRoot)
		buffer.WriteU32(hdr.Timestamp)
		buffer.WriteU32(hdr.Bits)
		buffer.WriteU32(hdr.Nonce)
		buffer.WriteU8(0x00)
	}
	return buffer.Bytes()
}

// Decode -- decoding from the bitcoin protocol format.
func (m *MsgHeaders) Decode(data []byte) error {
	var err error
	buffer := xbase.NewBufferReader(data)

	count, err := buffer.ReadVarInt()
	if err != nil {
		return err
	}

	for i := 0; i < int(count); i++ {
		header := &BlockHeader{}
		if header.Version, err = buffer.ReadU32(); err != nil {
			return err
		}
		if header.PrevBlock, err = buffer.ReadBytes(32); err != nil {
			return err
		}
		if header.MerkleRoot, err = buffer.ReadBytes(32); err != nil {
			return err
		}
		if header.Timestamp, err = buffer.ReadU32(); err != nil {
			return err
		}
		if header.Bits, err = buffer.ReadU32(); err != nil {
			return err
		}
		if header.Nonce, err = buffer.ReadU32(); err != nil {
			return err
		}
		if _, err := buffer.ReadU8(); err != nil {
			return err
		}
		m.Headers = append(m.Headers, header)
	}
	return nil
}

// Size -- the size of the message.
func (m *MsgHeaders) Size() int {
	size := xbase.VarIntSerializeSize(uint64(len(m.Headers)))
	for _, hdr := range m.Headers {
		size += hdr.Size() + 1
	}
	return size
}

// Command -- returns the protocol command of this message.
func (m *MsgHeaders) Command() string {
	return CommandHeaders
}
