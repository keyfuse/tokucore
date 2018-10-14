// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"github.com/tokublock/tokucore/network"
)

// *******************************************
// PayToWitnessPubKeyHashAddress(P2WPKH)
// *******************************************

// PayToWitnessPubKeyHashAddress -- is an Address for a pay-to-witness-pubkey-hash (P2WPKH) output.
// Public key -> P2WPKH address
// witness program = ripemd160(sha256(public key)).
// Encode into bech32 by providing the witness program, bc as the human readable part and 0 as witness version.
// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
type PayToWitnessPubKeyHashAddress struct {
	witnessVersion byte
	witnessProgram [20]byte
}

// NewPayToWitnessPubKeyHashAddress -- create new PayToWitnessPubKeyHashAddress.
func NewPayToWitnessPubKeyHashAddress(witnessProgram []byte) Address {
	if len(witnessProgram) != 20 {
		return nil
	}

	var witness [20]byte
	copy(witness[:], witnessProgram)
	return &PayToWitnessPubKeyHashAddress{
		witnessVersion: 0x00,
		witnessProgram: witness,
	}
}

// ToString -- the implementation method for xcore.Address interface.
func (a *PayToWitnessPubKeyHashAddress) ToString(net *network.Network) string {
	str, err := WitnessAddressEncode(net.Bech32HRPSegwit, a.witnessVersion, a.witnessProgram[:])
	if err != nil {
		return ""
	}
	return str
}

// Hash160 -- the address hash160 bytes.
func (a *PayToWitnessPubKeyHashAddress) Hash160() []byte {
	return a.witnessProgram[:]
}
