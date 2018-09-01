// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
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

// Script --
type Script interface {
	// GetAddress used to get the Address interface.
	GetAddress() Address
	// GetLockingScriptBytes used to get locking script bytes.
	GetLockingScriptBytes() ([]byte, error)
}

// isPubkeyHash --
// returns true if the script passed is a pay-to-pubkey-hash transaction, false otherwise.
func isPubkeyHash(instrs []xvm.Instruction) bool {
	return len(instrs) == 5 &&
		instrs[0].OpCode() == xvm.OP_DUP &&
		instrs[1].OpCode() == xvm.OP_HASH160 &&
		instrs[2].OpCode() == xvm.OP_DATA_20 &&
		instrs[3].OpCode() == xvm.OP_EQUALVERIFY &&
		instrs[4].OpCode() == xvm.OP_CHECKSIG
}

// isScriptHash --
// returns true if the script passed is a pay-to-script-hash transaction, false otherwise.
func isScriptHash(instrs []xvm.Instruction) bool {
	return len(instrs) == 3 &&
		instrs[0].OpCode() == xvm.OP_HASH160 &&
		instrs[1].OpCode() == xvm.OP_DATA_20 &&
		instrs[2].OpCode() == xvm.OP_EQUAL
}

func typeOfScript(instrs []xvm.Instruction) ScriptClass {
	switch {
	case isPubkeyHash(instrs):
		return PubKeyHashTy
	case isScriptHash(instrs):
		return ScriptHashTy
	}
	return NonStandardTy
}

// DataOutput -- returns OP_RETURN datas.
func DataOutput(data []byte) ([]byte, error) {
	return xvm.NewScriptBuilder().AddOp(xvm.OP_RETURN).AddData(data).Script()
}
