// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xerror"
	"github.com/tokublock/tokucore/xvm"
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
func NewPayToMultiSigScript(nrequired int, pkDatas ...[]byte) Script {
	return &PayToMultiSigScript{
		nrequired: nrequired,
		pkDatas:   pkDatas,
	}
}

// GetAddress -- returns the Address interface.
func (s *PayToMultiSigScript) GetAddress() Address {
	datas, err := s.GetLockingScriptBytes()
	if err != nil {
		return nil
	}
	return NewPayToScriptHashAddress(xcrypto.Hash160(datas))
}

// GetLockingScriptBytes -- returns the locking script bytes.
//
// OP_2 <A pubkey> <B pubkey> <C pubkey> OP_3 OP_CHECKMULTISIG
// Format:
// Required-N
// pubkey datas
// pubkey-N
// OP_CHECKMULTISIG
func (s *PayToMultiSigScript) GetLockingScriptBytes() ([]byte, error) {
	// Check the nrequired is valid
	if len(s.pkDatas) < s.nrequired {
		return nil, xerror.NewError(Errors, ER_SCRIPT_STANDARD_PUBKEYS_LE_NREQUIRED, len(s.pkDatas), s.nrequired)
	}
	builder := xvm.NewScriptBuilder()
	builder.AddInt64(int64(s.nrequired))
	for _, key := range s.pkDatas {
		builder.AddData(key)
	}
	builder.AddInt64(int64(len(s.pkDatas)))
	builder.AddOp(xvm.OP_CHECKMULTISIG)
	return builder.Script()
}
