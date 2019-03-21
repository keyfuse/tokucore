// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"github.com/tokublock/tokucore/xerror"
	"github.com/tokublock/tokucore/xvm"
)

// ScriptClass -- an enumeration for the list of standard types of script.
type ScriptClass byte

// Classes of script payment known about in the blockchain.
const (
	NonStandardTy         ScriptClass = iota // None of the recognized forms.
	PubKeyTy                                 // Pay pubkey.
	PubKeyHashTy                             // Pay pubkey hash.
	WitnessV0PubKeyHashTy                    // Pay witness pubkey hash.
	ScriptHashTy                             // Pay to script hash.
	WitnessV0ScriptHashTy                    // Pay to witness script hash.
	MultiSigTy                               // Multi signature.
	NullDataTy                               // Empty data-only (provably prunable).
)

// PubKeySign -- Public key and signature pair.
type PubKeySign struct {
	PubKey    []byte
	Signature []byte
}

// Script --
type Script interface {
	// GetAddress used to get the Address interface.
	GetAddress() Address

	// GetRawLockingScriptBytes -- used to get locking script bytes.
	GetRawLockingScriptBytes() ([]byte, error)

	// GetFinalLockingScriptBytes -- used to get the re-written locking for witness.
	// If txin is witness, returns re-written locking script code.
	// If txin is non-witness, returns same as GetLockingScriptBytes.
	GetFinalLockingScriptBytes(redeem []byte) ([]byte, error)

	// GetRawUnlockingScriptBytes -- used to get raw unlocking script bytes.
	GetRawUnlockingScriptBytes(signs []PubKeySign, redeem []byte) ([]byte, error)

	// GetWitnessUnlockingScriptBytes -- used to get witness script bytes.
	GetWitnessUnlockingScriptBytes(signs []PubKeySign, redeem []byte) ([][]byte, error)

	// GetWitnessScriptCode -- used to get the witness script for sighash of this txin.
	// If txin is non-witness, returns nil.
	GetWitnessScriptCode([]byte) ([]byte, error)

	// WitnessToUnlockingScriptBytes -- converts witness slice to unlocking script.
	// For txn deserialize from hex.
	WitnessToUnlockingScriptBytes(witness [][]byte) ([]byte, error)
}

// ParseLockingScript -- parse the locking script to script instance.
func ParseLockingScript(script []byte) (Script, error) {
	instrs, err := xvm.NewScriptReader(script).AllInstructions()
	if err != nil {
		return nil, err
	}
	switch {
	case isPubkeyHash(instrs):
		return NewPayToPubKeyHashScript(instrs[2].Data()), nil
	case isScriptHash(instrs):
		return NewPayToScriptHashScript(instrs[1].Data()), nil
	case isWitnessPubKeyHash(instrs):
		return NewPayToWitnessPubKeyHashScript(instrs[1].Data()), nil
	case isWitnessScriptHash(instrs):
		return NewPayToWitnessScriptHashScript(instrs[1].Data()), nil
	}
	return nil, xerror.NewError(Errors, ER_SCRIPT_TYPE_UNKNOWN, xvm.DisasmString(script))
}

// PayToAddrScript -- returns the locking script by address type.
func PayToAddrScript(addr Address) ([]byte, error) {
	switch addr.(type) {
	case *PayToPubKeyHashAddress:
		return NewPayToPubKeyHashScript(addr.Hash160()).GetRawLockingScriptBytes()
	case *PayToScriptHashAddress:
		return NewPayToScriptHashScript(addr.Hash160()).GetRawLockingScriptBytes()
	case *PayToWitnessPubKeyHashAddress:
		return NewPayToWitnessPubKeyHashScript(addr.Hash160()).GetRawLockingScriptBytes()
	case *PayToWitnessScriptHashAddress:
		return NewPayToWitnessScriptHashScript(addr.Hash160()).GetRawLockingScriptBytes()
	}
	return nil, xerror.NewError(Errors, ER_SCRIPT_STANDARD_ADDRESS_TYPE_UNSUPPORTED, addr)
}
