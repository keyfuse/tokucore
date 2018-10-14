// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"github.com/tokublock/tokucore/xvm"
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

// GetLockingScriptBytes -- returns the locking script bytes.
//
// OP_HASH160 <Hash160(redeemScript)> OP_EQUAL
// Format:
// - OP_HASH160
// - OP_DATA_20
// - 20 bytes script hash
// - OP_EQUAL
func (s *PayToScriptHashScript) GetLockingScriptBytes() ([]byte, error) {
	return xvm.NewScriptBuilder().
		AddOp(xvm.OP_HASH160).
		AddData(s.hash).
		AddOp(xvm.OP_EQUAL).
		Script()
}

// isScriptHash --
// returns true if the script passed is a pay-to-script-hash transaction, false otherwise.
func isScriptHash(instrs []xvm.Instruction) bool {
	return len(instrs) == 3 &&
		instrs[0].OpCode() == xvm.OP_HASH160 &&
		instrs[1].OpCode() == xvm.OP_DATA_20 &&
		instrs[2].OpCode() == xvm.OP_EQUAL
}
