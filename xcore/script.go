// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcore

import (
	"github.com/keyfuse/tokucore/xerror"
	"github.com/keyfuse/tokucore/xvm"
)

// ScriptVersion --
type ScriptVersion int

// Script version.
const (
	BASE ScriptVersion = iota
	WITNESS_V0
	TAPROOT
	TAPSCRIPT
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

	// GetScriptVersion -- used to get the version of the script.
	GetScriptVersion() ScriptVersion

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
	case isWitnessV0PubKeyHash(instrs):
		return NewPayToWitnessV0PubKeyHashScript(instrs[1].Data()), nil
	case isWitnessV0ScriptHash(instrs):
		return NewPayToWitnessV0ScriptHashScript(instrs[1].Data()), nil
	}
	return nil, xerror.NewError(Errors, ER_SCRIPT_TYPE_UNKNOWN, xvm.DisasmString(script))
}
