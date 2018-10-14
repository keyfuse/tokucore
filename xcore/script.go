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
	// GetLockingScriptBytes used to get locking script bytes.
	GetLockingScriptBytes() ([]byte, error)
}

// GetScriptClass -- gets the script class.
func GetScriptClass(script []byte) ScriptClass {
	instrs, err := xvm.NewScriptReader(script).AllInstructions()
	if err != nil {
		return NonStandardTy
	}
	switch {
	case isPubkeyHash(instrs):
		return PubKeyHashTy
	case isScriptHash(instrs):
		return ScriptHashTy
	case isWitnessPubKeyHash(instrs):
		return WitnessV0PubKeyHashTy
	}
	return NonStandardTy
}

// PayToAddrScript -- returns the locking script by address type.
func PayToAddrScript(addr Address) ([]byte, error) {
	switch addr.(type) {
	case *PayToPubKeyHashAddress:
		return NewPayToPubKeyHashScript(addr.Hash160()).GetLockingScriptBytes()
	case *PayToScriptHashAddress:
		return NewPayToScriptHashScript(addr.Hash160()).GetLockingScriptBytes()
	case *PayToWitnessPubKeyHashAddress:
		return NewPayToWitnessPubKeyHashScript(addr.Hash160()).GetLockingScriptBytes()
	}
	return nil, xerror.NewError(Errors, ER_SCRIPT_STANDARD_ADDRESS_TYPE_UNSUPPORTED, addr)
}
