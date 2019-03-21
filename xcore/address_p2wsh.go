// tokucore
//
// Copyright (c) 2018-2019 TokuBlock
// BSD License

package xcore

import (
	"github.com/tokublock/tokucore/network"
)

// *******************************************
// PayToWitnessScriptHashAddress(P2WSH)
// *******************************************

// PayToWitnessScriptHashAddress -- is an Address for a pay-to-witness-script-hash (P2WSH) output.
// witness program = sha256(script).
// Encode into bech32 by providing the witness program, bc as the human readable part and 0 as witness version.
// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
type PayToWitnessScriptHashAddress struct {
	witnessVersion byte
	witnessProgram [32]byte
}

// NewPayToWitnessScriptHashAddress -- create new PayToWitnessScriptHashAddress.
func NewPayToWitnessScriptHashAddress(witnessProgram []byte) Address {
	if len(witnessProgram) != 32 {
		return nil
	}

	var witness [32]byte
	copy(witness[:], witnessProgram)
	return &PayToWitnessScriptHashAddress{
		witnessVersion: 0x00,
		witnessProgram: witness,
	}
}

// ToString -- the implementation method for xcore.Address interface.
func (a *PayToWitnessScriptHashAddress) ToString(net *network.Network) string {
	str, err := WitnessAddressEncode(net.Bech32HRPSegwit, a.witnessVersion, a.witnessProgram[:])
	if err != nil {
		return ""
	}
	return str
}

// Hash160 -- the address hash160 bytes.
// Here is sha256, not hash160.
func (a *PayToWitnessScriptHashAddress) Hash160() []byte {
	return a.witnessProgram[:]
}
