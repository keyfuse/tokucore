// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcore

import (
	"github.com/keyfuse/tokucore/xvm"
)

// PayToScriptHashScript -- P2SH.
type PayToScriptHashScript struct {
	hash []byte
}

// NewPayToScriptHashScript -- creates P2SH.
func NewPayToScriptHashScript(scriptHash []byte) Script {
	return &PayToScriptHashScript{
		hash: scriptHash,
	}
}

// GetAddress -- returns the Address interface.
func (s *PayToScriptHashScript) GetAddress() Address {
	return NewPayToScriptHashAddress(s.hash)
}

// GetRawLockingScriptBytes -- used to get locking script bytes.
//
// OP_HASH160 <Hash160(redeemScript)> OP_EQUAL
// Format:
// - OP_HASH160
// - OP_DATA_20
// - 20 bytes script hash
// - OP_EQUAL
func (s *PayToScriptHashScript) GetRawLockingScriptBytes() ([]byte, error) {
	return xvm.NewScriptBuilder().
		AddOp(xvm.OP_HASH160).
		AddData(s.hash).
		AddOp(xvm.OP_EQUAL).
		Script()
}

// GetFinalLockingScriptBytes -- used to get the re-written locking for witness.
// Same as raw.
func (s *PayToScriptHashScript) GetFinalLockingScriptBytes(redeem []byte) ([]byte, error) {
	return s.GetRawLockingScriptBytes()
}

// GetRawUnlockingScriptBytes -- used to get raw unlocking script bytes.
// unlocking: OP_0 <A sig> <C sig> <redeemScript>
// witness:   (empty)
func (s *PayToScriptHashScript) GetRawUnlockingScriptBytes(signs []PubKeySign, redeem []byte) ([]byte, error) {
	builder := xvm.NewScriptBuilder()
	if len(signs) > 0 {
		builder.AddOp(xvm.OP_0)

		for _, sign := range signs {
			builder.AddData(sign.Signature)
		}
	}
	return builder.AddData(redeem).Script()
}

// GetWitnessUnlockingScriptBytes -- used to get witness script bytes.
func (s *PayToScriptHashScript) GetWitnessUnlockingScriptBytes(signs []PubKeySign, redeem []byte) ([][]byte, error) {
	return nil, nil
}

// GetWitnessScriptCode -- used to get the witness script for sighash of this txin.
func (s *PayToScriptHashScript) GetWitnessScriptCode(redeem []byte) ([]byte, error) {
	return nil, nil
}

// GetScriptVersion -- used to get the version of this script.
func (s *PayToScriptHashScript) GetScriptVersion() ScriptVersion {
	return BASE
}

// WitnessToUnlockingScriptBytes -- converts witness slice to unlocking script.
// For txn deserialize from hex.
func (s *PayToScriptHashScript) WitnessToUnlockingScriptBytes(witness [][]byte) ([]byte, error) {
	return nil, nil
}

// isScriptHash --
// returns true if the script passed is a pay-to-script-hash transaction, false otherwise.
func isScriptHash(instrs []xvm.Instruction) bool {
	return len(instrs) == 3 &&
		instrs[0].OpCode() == xvm.OP_HASH160 &&
		instrs[1].OpCode() == xvm.OP_DATA_20 &&
		instrs[2].OpCode() == xvm.OP_EQUAL
}
