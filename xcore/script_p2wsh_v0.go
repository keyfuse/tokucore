// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcore

import (
	"bytes"
	"fmt"

	"github.com/keyfuse/tokucore/xcrypto"
	"github.com/keyfuse/tokucore/xvm"
)

// PayToWitnessV0ScriptHashScript -- P2WSH (version 0 pay-to-witness-script-hash).
type PayToWitnessV0ScriptHashScript struct {
	hash []byte
}

// NewPayToWitnessV0ScriptHashScript -- creates new P2WSH script.
// scriptHash = sha256(script).
func NewPayToWitnessV0ScriptHashScript(scriptHash []byte) Script {
	return &PayToWitnessV0ScriptHashScript{
		hash: scriptHash,
	}
}

// GetAddress -- returns the Address interface.
func (s *PayToWitnessV0ScriptHashScript) GetAddress() Address {
	return NewPayToWitnessV0ScriptHashAddress(s.hash)
}

// GetRawLockingScriptBytes -- used to get locking script bytes.
//
// 0 <32-byte-script-hash>
// Format:
// - OP_0
// - OP_DATA_32
// - 32 bytes sha256<script-hash>
func (s *PayToWitnessV0ScriptHashScript) GetRawLockingScriptBytes() ([]byte, error) {
	return xvm.NewScriptBuilder().
		AddOp(xvm.OP_0).
		AddData(s.hash).
		Script()
}

// GetFinalLockingScriptBytes -- used to get the re-written locking for witness.
func (s *PayToWitnessV0ScriptHashScript) GetFinalLockingScriptBytes(redeem []byte) ([]byte, error) {
	scriptInstance := NewPayToScriptHashScript(xcrypto.Hash160(redeem))
	return scriptInstance.GetFinalLockingScriptBytes(redeem)
}

// GetRawUnlockingScriptBytes -- used to get raw unlocking script bytes.
// unlocking: (empty)
func (s *PayToWitnessV0ScriptHashScript) GetRawUnlockingScriptBytes(signs []PubKeySign, redeem []byte) ([]byte, error) {
	// Check sha256(redeem) == s.Hash
	if !bytes.Equal(xcrypto.Sha256(redeem), s.hash) {
		return nil, fmt.Errorf("PayToWitnessV0ScriptHashScript.GetUnlockingScriptBytes.error:sha256(redeem)!=s.hash")
	}

	scriptInstance := NewPayToScriptHashScript(xcrypto.Hash160(redeem))
	return scriptInstance.GetRawUnlockingScriptBytes(signs, redeem)
}

// GetWitnessUnlockingScriptBytes -- used to get witness script bytes.
func (s *PayToWitnessV0ScriptHashScript) GetWitnessUnlockingScriptBytes(signs []PubKeySign, redeem []byte) ([][]byte, error) {
	// Check sha256(redeem) == s.Hash
	if !bytes.Equal(xcrypto.Sha256(redeem), s.hash) {
		return nil, fmt.Errorf("PayToWitnessV0ScriptHashScript.GetUnlockingScriptBytes.error:sha256(redeem)!=s.hash")
	}

	var witness [][]byte
	if len(signs) > 0 {
		// Dummy for CHECKMULTSIG.
		witness = append(witness, []byte{})
		for _, sign := range signs {
			witness = append(witness, sign.Signature)
		}
	}
	witness = append(witness, redeem)
	return witness, nil
}

// GetWitnessScriptCode -- used to get the witness script for sighash of this txin.
// For P2WSH witness, we convert it to P2SH.
// OP_0 OP_DATA_20 <20-bytes-script-hash>
// to
// OP_HASH160 <Hash160(redeemScript)> OP_EQUAL
func (s *PayToWitnessV0ScriptHashScript) GetWitnessScriptCode(redeem []byte) ([]byte, error) {
	return redeem, nil
}

// GetScriptVersion -- used to get the version of this script.
func (s *PayToWitnessV0ScriptHashScript) GetScriptVersion() ScriptVersion {
	return WITNESS_V0
}

// WitnessToUnlockingScriptBytes -- converts witness slice to unlocking script.
// For txn deserialize from hex.
func (s *PayToWitnessV0ScriptHashScript) WitnessToUnlockingScriptBytes(witness [][]byte) ([]byte, error) {
	l := len(witness)
	if l > 1 {
		redeem := witness[len(witness)-1]
		instrs, _ := xvm.NewScriptReader(redeem).AllInstructions()
		if isMultiSig(instrs) {
			builder := xvm.NewScriptBuilder()
			for i := 0; i < l-1; i++ {
				builder.AddData(witness[i])
			}
			return builder.AddData(redeem).Script()
		}
	}
	return nil, nil
}

// isWitnessV0ScriptHash --
// returns true if the passed script is a pay-to-witness-script-hash, and false otherwise.
func isWitnessV0ScriptHash(instrs []xvm.Instruction) bool {
	return len(instrs) == 2 &&
		instrs[0].OpCode() == xvm.OP_0 &&
		instrs[1].OpCode() == xvm.OP_DATA_32
}
