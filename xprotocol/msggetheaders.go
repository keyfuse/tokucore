// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xbase"
)

// MsgGetHeaders --
type MsgGetHeaders struct {
	protocolVersion    uint32
	blockLocatorHashes [][]byte
	blockStopHash      [32]byte
}

// NewMsgGetHeaders -- creates new MsgGetHeaders.
func NewMsgGetHeaders(net *network.Network) *MsgGetHeaders {
	return &MsgGetHeaders{
		protocolVersion: net.ProtocolVersion,
	}
}

// AddBlockLocatorHash -- adds the block locator hash to message.
func (m *MsgGetHeaders) AddBlockLocatorHash(hash []byte) error {
	m.blockLocatorHashes = append(m.blockLocatorHashes, hash)
	return nil
}

// Encode -- encoding MsgGetHeaders.
func (m *MsgGetHeaders) Encode() []byte {
	if len(m.blockLocatorHashes) == 0 {
		var zero [32]byte
		m.AddBlockLocatorHash(zero[:])
	}

	buffer := xbase.NewBuffer()
	buffer.WriteU32(m.protocolVersion)
	buffer.WriteVarInt(uint64(len(m.blockLocatorHashes)))
	for _, locator := range m.blockLocatorHashes {
		buffer.WriteBytes(locator[:])
	}
	buffer.WriteBytes(m.blockStopHash[:])
	return buffer.Bytes()
}

// Decode -- decoding bytes to MsgGetHeaders.
func (m *MsgGetHeaders) Decode(data []byte) error {
	var err error
	buffer := xbase.NewBufferReader(data)

	if m.protocolVersion, err = buffer.ReadU32(); err != nil {
		return err
	}
	count, err := buffer.ReadVarInt()
	if err != nil {
		return err
	}
	for i := 0; i < int(count); i++ {
		var hash []byte
		if hash, err = buffer.ReadBytes(32); err != nil {
			return err
		}
		m.AddBlockLocatorHash(hash)
	}
	return nil
}

// Size -- the size of the message.
func (m *MsgGetHeaders) Size() int {
	size := 4
	size += xbase.VarIntSerializeSize(uint64(len(m.blockLocatorHashes)))
	for _, hdr := range m.blockLocatorHashes {
		size += len(hdr)
	}
	size += len(m.blockStopHash)
	return size
}

// Command -- returns the protocol command of this message.
func (m *MsgGetHeaders) Command() string {
	return CommandGetHeaders
}
