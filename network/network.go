// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package network

// Network -- network paramars.
type Network struct {
	// Address encoding magics.
	PubKeyHashAddrID        byte // First byte of a P2PKH address
	ScriptHashAddrID        byte // First byte of a P2PSH address
	PrivateKeyID            byte // First byte of a WIF private key
	WitnessPubKeyHashAddrID byte // First byte of a P2WPKH address
	WitnessScriptHashAddrID byte // First byte of P2WSH address

	// Human-readable part for Bech32 encoded segwit addresses, as defined in BIP173.
	Bech32HRPSegwit string

	// BIP32 hierarchical deterministic extended key magics.
	HDPrivateKeyID []byte
	HDPublicKeyID  []byte

	// Protocol.
	Magic           []byte
	Port            uint32
	LastBlock       uint32
	ProtocolVersion uint32
	UserAgent       string
}

var (
	// MainNet -- bitcoin mainet.
	MainNet = &Network{
		// Address encoding magics.
		PubKeyHashAddrID:        0x00, // starts with 1
		ScriptHashAddrID:        0x05, // starts with 3
		PrivateKeyID:            0x80, // starts with 5 (uncompressed) or K (compressed)
		WitnessPubKeyHashAddrID: 0x06, // starts with p2
		WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

		// Human-readable part for Bech32 encoded segwit addresses, as defined in BIP173.
		Bech32HRPSegwit: "bc", // always bc for main net

		// BIP32 hierarchical deterministic extended key magics
		HDPrivateKeyID: []byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
		HDPublicKeyID:  []byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

		// Protocol.
		Magic:           []byte{0xf9, 0xbe, 0xb4, 0xd9},
		Port:            8333,
		LastBlock:       1,
		ProtocolVersion: 70015,
		UserAgent:       "/tokucore:0.0.1/",
	}

	// TestNet -- bitcoin testnet.
	TestNet = &Network{
		// Address encoding magics.
		PubKeyHashAddrID:        0x6f, // starts with m or n
		ScriptHashAddrID:        0xc4, // starts with 2
		WitnessPubKeyHashAddrID: 0x03, // starts with QW
		WitnessScriptHashAddrID: 0x28, // starts with T7n
		PrivateKeyID:            0xef, // starts with 9 (uncompressed) or c (compressed)

		// Human-readable part for Bech32 encoded segwit addresses, as defined in BIP173.
		Bech32HRPSegwit: "tb", // always tb for test net

		// BIP32 hierarchical deterministic extended key magics
		HDPrivateKeyID: []byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
		HDPublicKeyID:  []byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

		// Protocol.
		Magic:           []byte{0x0b, 0x11, 0x09, 0x07},
		Port:            18333,
		LastBlock:       1,
		ProtocolVersion: 70015,
		UserAgent:       "/tokucore:0.0.1/",
	}
)

// SetMagic -- set magic.
func (net *Network) SetMagic(magic []byte) *Network {
	net.Magic = magic
	return net
}

// SetPort -- set port.
func (net *Network) SetPort(port uint32) *Network {
	net.Port = port
	return net
}

// SetLastBlock -- set last block.
func (net *Network) SetLastBlock(last uint32) *Network {
	net.LastBlock = last
	return net
}
