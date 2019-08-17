// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcore

import (
	"github.com/keyfuse/tokucore/xvm"
)

// PayToPubKeyHashScript -- P2PKH.
type PayToPubKeyHashScript struct {
	hash []byte
}

// NewPayToPubKeyHashScript -- creates new P2PKH.
func NewPayToPubKeyHashScript(pubkeyHash []byte) Script {
	return &PayToPubKeyHashScript{
		hash: pubkeyHash,
	}
}

// GetAddress -- returns the Address interface.
func (s *PayToPubKeyHashScript) GetAddress() Address {
	return NewPayToPubKeyHashAddress(s.hash)
}

// GetRawLockingScriptBytes -- used to get locking script bytes.
//
// OP_DUP OP_HASH160 <PubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
// Format:
// - OP_DUP
// - OP_HASH160
// - OP_DATA_20
// - 20 bytes pubkey hash
// - OP_EQUALVERIFY
// - OP_CHECKSIG
func (s *PayToPubKeyHashScript) GetRawLockingScriptBytes() ([]byte, error) {
	return xvm.NewScriptBuilder().
		AddOp(xvm.OP_DUP).
		AddOp(xvm.OP_HASH160).
		AddData(s.hash).
		AddOp(xvm.OP_EQUALVERIFY).
		AddOp(xvm.OP_CHECKSIG).
		Script()
}

// GetFinalLockingScriptBytes -- used to get the re-written locking for witness.
// Same as raw.
func (s *PayToPubKeyHashScript) GetFinalLockingScriptBytes(redeem []byte) ([]byte, error) {
	return s.GetRawLockingScriptBytes()
}

// GetRawUnlockingScriptBytes -- returns the unlocking script bytes.
// unlocking: <sig> <pubkey>
// witness:   (empty)
func (s *PayToPubKeyHashScript) GetRawUnlockingScriptBytes(signs []PubKeySign, redeem []byte) ([]byte, error) {
	builder := xvm.NewScriptBuilder()
	unlocking, err := builder.AddData(signs[0].Signature).AddData(signs[0].PubKey).Script()
	return unlocking, err
}

// GetWitnessUnlockingScriptBytes -- used to get witness script bytes.
func (s *PayToPubKeyHashScript) GetWitnessUnlockingScriptBytes(signs []PubKeySign, redeem []byte) ([][]byte, error) {
	return nil, nil
}

// GetWitnessScriptCode -- used to get the unlocking script bytes of witness program.
func (s *PayToPubKeyHashScript) GetWitnessScriptCode(redeem []byte) ([]byte, error) {
	return nil, nil
}

// GetScriptVersion -- used to get the version of this script.
func (s *PayToPubKeyHashScript) GetScriptVersion() ScriptVersion {
	return BASE
}

// WitnessToUnlockingScriptBytes -- converts witness slice to unlocking script.
// For txn deserialize from hex.
func (s *PayToPubKeyHashScript) WitnessToUnlockingScriptBytes(witness [][]byte) ([]byte, error) {
	return nil, nil
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
