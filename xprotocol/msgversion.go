// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"time"

	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xbase"
)

const (
	ip = "::ffff:127.0.0.1" // always return loopback.
)

// MsgVersion --
type MsgVersion struct {
	// Version of the protocol the node is using.
	Version uint32

	// Bitfield which identifies the enabled services.
	// 0: no services supported on this node.
	Services uint64

	// Time the message was generated.
	Timestamp uint64

	// Services  of the remote peer.
	// Bitfield which identifies the services supported by the address.
	ServicesYou uint64

	// Address of the remote peer.
	AddressYou []byte

	// Port of the remote peer.
	PortYou uint32

	// Services  of the local peer.
	// Bitfield which identifies the services supported by the address.
	ServicesMe uint64

	// Address of the local peer.
	AddressMe []byte

	// Port of the local peer.
	PortMe uint32

	// Unique value associated with message that is used to detect self
	// connections.
	Nonce uint64

	// The user agent that generated messsage.
	// This is a encoded as a varString.
	UserAgent string

	// Last block seen by the generator of the version message.
	LastBlock uint32

	// Announce transactions to peer.
	Relay byte
}

// NewMsgVersion -- creates new MsgVersion.
func NewMsgVersion(network *network.Network) *MsgVersion {
	return &MsgVersion{
		Version:    network.ProtocolVersion,
		Timestamp:  uint64(time.Now().Unix()),
		AddressYou: []byte(ip),
		PortYou:    network.Port,
		AddressMe:  []byte(ip),
		PortMe:     network.Port,
		UserAgent:  network.UserAgent,
		LastBlock:  1,
	}
}

// Encode -- encoding to bitcoin protocol format.
func (m *MsgVersion) Encode() []byte {
	buffer := xbase.NewBuffer()
	buffer.WriteU32((m.Version))
	buffer.WriteU64((m.Services))
	buffer.WriteU64(m.Timestamp)
	buffer.WriteU64(m.ServicesYou)
	buffer.WriteBytes(m.AddressYou)
	buffer.WriteU16(m.PortYou)
	buffer.WriteU64(m.ServicesMe)
	buffer.WriteBytes(m.AddressMe)
	buffer.WriteU16(m.PortMe)
	buffer.WriteU64(m.Nonce)
	buffer.WriteVarString(m.UserAgent)
	buffer.WriteU32((m.LastBlock))
	buffer.WriteU8(m.Relay)
	return buffer.Bytes()
}

// Decode -- decoding from the bitcoin protocol format.
func (m *MsgVersion) Decode(data []byte) error {
	var err error
	buffer := xbase.NewBufferReader(data)

	if m.Version, err = buffer.ReadU32(); err != nil {
		return err
	}
	if m.Services, err = buffer.ReadU64(); err != nil {
		return err
	}
	if m.Timestamp, err = buffer.ReadU64(); err != nil {
		return err
	}
	if m.ServicesYou, err = buffer.ReadU64(); err != nil {
		return err
	}
	if m.AddressYou, err = buffer.ReadBytes(16); err != nil {
		return err
	}
	if m.PortYou, err = buffer.ReadU16(); err != nil {
		return err
	}
	if m.ServicesMe, err = buffer.ReadU64(); err != nil {
		return err
	}
	if m.AddressMe, err = buffer.ReadBytes(16); err != nil {
		return err
	}
	if m.PortMe, err = buffer.ReadU16(); err != nil {
		return err
	}
	if m.Nonce, err = buffer.ReadU64(); err != nil {
		return err
	}
	if m.UserAgent, err = buffer.ReadVarString(); err != nil {
		return err
	}
	if m.LastBlock, err = buffer.ReadU32(); err != nil {
		return err
	}
	if m.Relay, err = buffer.ReadU8(); err != nil {
		return err
	}
	return nil
}

// Size -- the size of the message.
func (m *MsgVersion) Size() int {
	return 4 + 8 + 8 + 8 +
		len(m.AddressYou) + 2 + 8 +
		len(m.AddressMe) + 2 + 8 +
		xbase.VarIntSerializeSize(uint64(len(m.UserAgent))) + len(m.UserAgent) + 4 + 1
}

// Command -- returns the protocol command of this message.
func (m *MsgVersion) Command() string {
	return CommandVersion
}
