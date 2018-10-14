// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"github.com/tokublock/tokucore/xvm"
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

// GetLockingScriptBytes --
// returns the locking script bytes.
//
// OP_DUP OP_HASH160 <PubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
// Format:
// - OP_DUP
// - OP_HASH160
// - OP_DATA_20
// - 20 bytes pubkey hash
// - OP_EQUALVERIFY
// - OP_CHECKSIG
func (s *PayToPubKeyHashScript) GetLockingScriptBytes() ([]byte, error) {
	return xvm.NewScriptBuilder().
		AddOp(xvm.OP_DUP).
		AddOp(xvm.OP_HASH160).
		AddData(s.hash).
		AddOp(xvm.OP_EQUALVERIFY).
		AddOp(xvm.OP_CHECKSIG).
		Script()
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
