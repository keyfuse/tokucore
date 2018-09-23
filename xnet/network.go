// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xnet

// Network --
type Network struct {
	Magic           []byte
	Port            uint32
	LastBlock       uint32
	ProtocolVersion uint32
	UserAgent       string
}

var (
	// MainNet --
	MainNet = &Network{
		Magic:           []byte{0xf9, 0xbe, 0xb4, 0xd9},
		Port:            8333,
		LastBlock:       1,
		ProtocolVersion: 70015,
		UserAgent:       "/tokucore:0.0.1/",
	}

	// TestNet --
	TestNet = &Network{
		Magic:           []byte{0x0b, 0x11, 0x09, 0x07},
		Port:            18333,
		LastBlock:       1,
		ProtocolVersion: 70015,
		UserAgent:       "/tokucore:0.0.1/",
	}
)
