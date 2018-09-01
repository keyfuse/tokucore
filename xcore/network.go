// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

// Network --
type Network interface {
	CreatePubkeyAddress([]byte) (Address, error)
	CreatePubkeyHashAddress([]byte) (Address, error)
}

// NetworkParams -- network paramars.
type NetworkParams struct {
	// Address encoding magics.
	PubKeyHashAddrID        byte // First byte of a P2PKH address
	ScriptHashAddrID        byte // First byte of a P2PSH address
	PrivateKeyID            byte // First byte of a WIF private key
	WitnessPubKeyHashAddrID byte // First byte of a P2WPKH address
	WitnessScriptHashAddrID byte // First byte of P2WSH address

	// BIP32 hierarchical deterministic extended key magics.
	HDPrivateKeyID []byte
	HDPublicKeyID  []byte
}

var (
	// MainNet -- bitcoin mainet.
	MainNet = &NetworkParams{
		// Address encoding magics.
		PubKeyHashAddrID:        0x00, // starts with 1
		ScriptHashAddrID:        0x05, // starts with 3
		PrivateKeyID:            0x80, // starts with 5 (uncompressed) or K (compressed)
		WitnessPubKeyHashAddrID: 0x06, // starts with p2
		WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

		// BIP32 hierarchical deterministic extended key magics
		HDPrivateKeyID: []byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
		HDPublicKeyID:  []byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub
	}

	// TestNet -- bitcoin testnet.
	TestNet = &NetworkParams{
		// Address encoding magics.
		PubKeyHashAddrID:        0x6f, // starts with m or n
		ScriptHashAddrID:        0xc4, // starts with 2
		WitnessPubKeyHashAddrID: 0x03, // starts with QW
		WitnessScriptHashAddrID: 0x28, // starts with T7n
		PrivateKeyID:            0xef, // starts with 9 (uncompressed) or c (compressed)

		// BIP32 hierarchical deterministic extended key magics
		HDPrivateKeyID: []byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
		HDPublicKeyID:  []byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub
	}
)
