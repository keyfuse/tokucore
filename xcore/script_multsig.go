// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcore

import (
	"github.com/keyfuse/tokucore/xcrypto"
	"github.com/keyfuse/tokucore/xerror"
	"github.com/keyfuse/tokucore/xvm"
)

// PayToMultiSigScript --
// returns a valid script for a multi-signature redemption where n-required of
// keys in pubkeys are required to have signed the transaction for success.
// pubkey is the public-key with compressed/uncompress format.
type PayToMultiSigScript struct {
	nrequired int
	pkDatas   [][]byte
}

// NewPayToMultiSigScript -- creates new PayToMultiSigScript.
func NewPayToMultiSigScript(nrequired int, pkDatas ...[]byte) *PayToMultiSigScript {
	return &PayToMultiSigScript{
		nrequired: nrequired,
		pkDatas:   pkDatas,
	}
}

// GetLockingScriptBytes -- returns the locking script bytes.
func (s *PayToMultiSigScript) GetLockingScriptBytes() ([]byte, error) {
	return GenMultiSigScript(s.nrequired, s.pkDatas...)
}

// GenMultiSigScript -- generates the multisig script for N of M.
//
// OP_2 <A pubkey> <B pubkey> <C pubkey> OP_3 OP_CHECKMULTISIG
// Format:
// Required-N
// pubkey datas
// pubkey-N
// OP_CHECKMULTISIG
func GenMultiSigScript(nrequired int, pkDatas ...[]byte) ([]byte, error) {
	// Check the nrequired is valid
	if len(pkDatas) < nrequired {
		return nil, xerror.NewError(Errors, ER_SCRIPT_STANDARD_PUBKEYS_LE_NREQUIRED, len(pkDatas), nrequired)
	}
	builder := xvm.NewScriptBuilder()
	builder.AddInt64(int64(nrequired))
	for _, key := range pkDatas {
		builder.AddData(key)
	}
	builder.AddInt64(int64(len(pkDatas)))
	builder.AddOp(xvm.OP_CHECKMULTISIG)
	return builder.Script()
}

// Hash160 -- returns the hash160(locking-script) bytes.
func (s *PayToMultiSigScript) Hash160() ([]byte, error) {
	datas, err := GenMultiSigScript(s.nrequired, s.pkDatas...)
	if err != nil {
		return nil, err
	}
	return xcrypto.Hash160(datas), nil
}

// Sha256 -- returns the sha256(locking-script) bytes.
func (s *PayToMultiSigScript) Sha256() ([]byte, error) {
	datas, err := GenMultiSigScript(s.nrequired, s.pkDatas...)
	if err != nil {
		return nil, err
	}
	return xcrypto.Sha256(datas), nil
}

// isSmallInt returns whether or not the opcode is considered a small integer,
// which is an OP_0, or OP_1 through OP_16.
func isSmallInt(op *xvm.Instruction) bool {
	if op.OpCode() == xvm.OP_0 || (op.OpCode() >= xvm.OP_1 && op.OpCode() <= xvm.OP_16) {
		return true
	}
	return false
}

// asSmallInt returns the passed opcode, which must be true according to
// isSmallInt(), as an integer.
func asSmallInt(op *xvm.Instruction) int {
	if op.OpCode() == xvm.OP_0 {
		return 0
	}
	return int(op.OpCode() - (xvm.OP_1 - 1))
}

func isMultiSig(instrs []xvm.Instruction) bool {
	l := len(instrs)

	// The absolute minimum is 1 pubkey:
	// OP_0/OP_1-16 <pubkey> OP_1 OP_CHECKMULTISIG
	if l < 4 {
		return false
	}

	if !isSmallInt(&instrs[0]) {
		return false
	}
	if !isSmallInt(&instrs[l-2]) {
		return false
	}
	if instrs[l-1].OpCode() != xvm.OP_CHECKMULTISIG {
		return false
	}

	// Verify the number of pubkeys specified matches the actual number
	// of pubkeys provided.
	if l-2-1 != asSmallInt(&instrs[l-2]) {
		return false
	}

	for _, instr := range instrs[1 : l-2] {
		// Valid pubkeys are either 33 or 65 bytes.
		if len(instr.Data()) != 33 && len(instr.Data()) != 65 {
			return false
		}
	}
	return true
}
