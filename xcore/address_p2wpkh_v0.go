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
// PayToWitnessV0PubKeyHashAddress(P2WPKH)
// *******************************************

// PayToWitnessV0PubKeyHashAddress -- is an Address for a pay-to-witness-pubkey-hash (P2WPKH) output.
// Public key -> P2WPKH address
// witness program = ripemd160(sha256(public key)).
// Encode into bech32 by providing the witness program, bc as the human readable part and 0 as witness version.
// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
type PayToWitnessV0PubKeyHashAddress struct {
	witnessVersion byte
	witnessProgram [20]byte
}

// NewPayToWitnessV0PubKeyHashAddress -- create new PayToWitnessV0PubKeyHashAddress.
func NewPayToWitnessV0PubKeyHashAddress(witnessProgram []byte) Address {
	if len(witnessProgram) != 20 {
		return nil
	}

	var witness [20]byte
	copy(witness[:], witnessProgram)
	return &PayToWitnessV0PubKeyHashAddress{
		witnessVersion: 0x00,
		witnessProgram: witness,
	}
}

// ToString -- the implementation method for xcore.Address interface.
func (a *PayToWitnessV0PubKeyHashAddress) ToString(net *network.Network) string {
	str, err := xbase.WitnessEncode(net.Bech32HRPSegwit, a.witnessVersion, a.witnessProgram[:])
	if err != nil {
		return ""
	}
	return str
}

// Hash160 -- the address hash160 bytes.
func (a *PayToWitnessV0PubKeyHashAddress) Hash160() []byte {
	return a.witnessProgram[:]
}

// LockingScript -- the address locking script.
func (a *PayToWitnessV0PubKeyHashAddress) LockingScript() ([]byte, error) {
	return NewPayToWitnessV0PubKeyHashScript(a.Hash160()).GetRawLockingScriptBytes()
}
