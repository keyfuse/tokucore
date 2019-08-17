// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcore

import (
	"github.com/keyfuse/tokucore/xvm"
)

// PayToWitnessV0PubKeyHashScript -- P2WPKH (version 0 pay-to-witness-public-key-hash).
type PayToWitnessV0PubKeyHashScript struct {
	hash []byte
}

// NewPayToWitnessV0PubKeyHashScript -- creates new P2WPKH script.
// pubkeyHash = hash160(pubkey)
func NewPayToWitnessV0PubKeyHashScript(pubkeyHash []byte) Script {
	return &PayToWitnessV0PubKeyHashScript{
		hash: pubkeyHash,
	}
}

// GetAddress -- returns the Address interface.
func (s *PayToWitnessV0PubKeyHashScript) GetAddress() Address {
	return NewPayToWitnessV0PubKeyHashAddress(s.hash)
}

// GetRawLockingScriptBytes -- used to get locking script bytes.
//
// 0 <20-byte-key-hash>
// Format:
// - OP_0
// - OP_DATA_20
// - 20 bytes pubkey hash
func (s *PayToWitnessV0PubKeyHashScript) GetRawLockingScriptBytes() ([]byte, error) {
	return xvm.NewScriptBuilder().
		AddOp(xvm.OP_0).
		AddData(s.hash).
		Script()
}

// GetFinalLockingScriptBytes -- used to get the re-written locking for witness.
// Same as P2PKH.
func (s *PayToWitnessV0PubKeyHashScript) GetFinalLockingScriptBytes(redeem []byte) ([]byte, error) {
	scriptInstance := NewPayToPubKeyHashScript(s.hash)
	return scriptInstance.GetRawLockingScriptBytes()
}

// GetRawUnlockingScriptBytes -- used to get raw unlocking script bytes.
func (s *PayToWitnessV0PubKeyHashScript) GetRawUnlockingScriptBytes(signs []PubKeySign, redeem []byte) ([]byte, error) {
	builder := xvm.NewScriptBuilder()
	return builder.AddData(signs[0].Signature).AddData(signs[0].PubKey).Script()
}

// GetWitnessUnlockingScriptBytes -- used to get witness script bytes.
func (s *PayToWitnessV0PubKeyHashScript) GetWitnessUnlockingScriptBytes(signs []PubKeySign, redeem []byte) ([][]byte, error) {
	var witness [][]byte
	witness = append(witness, signs[0].Signature)
	witness = append(witness, signs[0].PubKey)
	return witness, nil
}

// GetWitnessScriptCode -- used to get the normal locking script bytes of witness program.
// For P2WPKH witness, we convert it to P2PKH.
// OP_0 OP_DATA_20 <20-bytes-pubkey-hash>
// to
// OP_DUP OP_HASH160 <20-bytes-pubkey-hash> OP_EQUALVERIFY OP_CHECKSIG
func (s *PayToWitnessV0PubKeyHashScript) GetWitnessScriptCode(redeem []byte) ([]byte, error) {
	scriptInstance := NewPayToPubKeyHashScript(s.hash)
	return scriptInstance.GetRawLockingScriptBytes()
}

// GetScriptVersion -- used to get the version of this script.
func (s *PayToWitnessV0PubKeyHashScript) GetScriptVersion() ScriptVersion {
	return WITNESS_V0
}

// WitnessToUnlockingScriptBytes -- converts witness slice to unlocking script.
// For txn deserialize from hex.
func (s *PayToWitnessV0PubKeyHashScript) WitnessToUnlockingScriptBytes(witness [][]byte) ([]byte, error) {
	builder := xvm.NewScriptBuilder()
	return builder.AddData(witness[0]).AddData(witness[1]).Script()
}

// isWitnessV0PubKeyHash --
// returns true if the passed script is a pay-to-witness-pubkey-hash, and false otherwise.
func isWitnessV0PubKeyHash(instrs []xvm.Instruction) bool {
	return len(instrs) == 2 &&
		instrs[0].OpCode() == xvm.OP_0 &&
		instrs[1].OpCode() == xvm.OP_DATA_20
}
