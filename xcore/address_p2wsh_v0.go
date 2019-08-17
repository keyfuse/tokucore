// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcore

import (
	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xbase"
)

// *******************************************
// PayToWitnessV0ScriptHashAddress(P2WSH)
// *******************************************

// PayToWitnessV0ScriptHashAddress -- is an Address for a pay-to-witness-script-hash (P2WSH) output.
// witness program = sha256(script).
// Encode into bech32 by providing the witness program, bc as the human readable part and 0 as witness version.
// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
type PayToWitnessV0ScriptHashAddress struct {
	witnessVersion byte
	witnessProgram [32]byte
}

// NewPayToWitnessV0ScriptHashAddress -- create new PayToWitnessV0ScriptHashAddress.
func NewPayToWitnessV0ScriptHashAddress(witnessProgram []byte) Address {
	if len(witnessProgram) != 32 {
		return nil
	}

	var witness [32]byte
	copy(witness[:], witnessProgram)
	return &PayToWitnessV0ScriptHashAddress{
		witnessVersion: 0x00,
		witnessProgram: witness,
	}
}

// ToString -- the implementation method for xcore.Address interface.
func (a *PayToWitnessV0ScriptHashAddress) ToString(net *network.Network) string {
	str, err := xbase.WitnessEncode(net.Bech32HRPSegwit, a.witnessVersion, a.witnessProgram[:])
	if err != nil {
		return ""
	}
	return str
}

// Hash160 -- the address hash160 bytes.
// Here is sha256, not hash160.
func (a *PayToWitnessV0ScriptHashAddress) Hash160() []byte {
	return a.witnessProgram[:]
}

// LockingScript -- the address locking script.
func (a *PayToWitnessV0ScriptHashAddress) LockingScript() ([]byte, error) {
	return NewPayToWitnessV0ScriptHashScript(a.Hash160()).GetRawLockingScriptBytes()
}
