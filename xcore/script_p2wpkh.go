// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"github.com/tokublock/tokucore/xvm"
)

// PayToWitnessPubKeyHashScript -- P2WPKH (version 0 pay-to-witness-public-key-hash).
type PayToWitnessPubKeyHashScript struct {
	hash []byte
}

// NewPayToWitnessPubKeyHashScript -- creates new P2WPKH.
func NewPayToWitnessPubKeyHashScript(pubkeyHash []byte) Script {
	return &PayToWitnessPubKeyHashScript{
		hash: pubkeyHash,
	}
}

// GetAddress -- returns the Address interface.
func (s *PayToWitnessPubKeyHashScript) GetAddress() Address {
	return NewPayToWitnessPubKeyHashAddress(s.hash)
}

// GetLockingScriptBytes --
// returns the locking script bytes.
//
// 0 <20-byte-key-hash>
// Format:
// - OP_0
// - OP_DATA_20
// - 20 bytes pubkey hash
func (s *PayToWitnessPubKeyHashScript) GetLockingScriptBytes() ([]byte, error) {
	return xvm.NewScriptBuilder().
		AddOp(xvm.OP_0).
		AddData(s.hash).
		Script()
}

// isWitnessPubKeyHash --
// returns true if the passed script is a pay-to-witness-pubkey-hash, and false otherwise.
func isWitnessPubKeyHash(instrs []xvm.Instruction) bool {
	return len(instrs) == 2 &&
		instrs[0].OpCode() == xvm.OP_0 &&
		instrs[1].OpCode() == xvm.OP_DATA_20
}
