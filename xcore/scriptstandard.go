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
// OP_DUP OP_HASH160 <hash> OP_EQUALVERIFY OP_CHECKSIG
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
// OP_HASH160 <scripthash> OP_EQUAL
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

// PayToAddrScript -- returns the locking script by address type.
func PayToAddrScript(addr Address) ([]byte, error) {
	switch addr.(type) {
	case *PayToPubKeyHashAddress:
		return NewPayToPubKeyHashScript(addr.Hash160()).GetLockingScriptBytes()
	case *PayToScriptHashAddress:
		return NewPayToScriptHashScript(addr.Hash160()).GetLockingScriptBytes()
	}
	return nil, xerror.NewError(Errors, ER_SCRIPT_STANDARD_ADDRESS_TYPE_UNSUPPORTED, addr)
}

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

// PubKeySign -- Public key and signature pair.
type PubKeySign struct {
	PubKey    []byte
	Signature []byte
}

// BuildUnlockingScriptBytes -- build the unlocking script bytes with redeemscript, pubkey and signature.
func BuildUnlockingScriptBytes(script []byte, redeemScript []byte, signs []PubKeySign) ([]byte, error) {
	pops, err := xvm.NewScriptReader(script).AllInstructions()
	if err != nil {
		return nil, err
	}

	builder := xvm.NewScriptBuilder()
	class := typeOfScript(pops)
	switch class {
	case PubKeyHashTy:
		sign := signs[0]
		// Format
		// Sig | pubkey
		return builder.AddData(sign.Signature).AddData(sign.PubKey).Script()
	case ScriptHashTy:
		// Format
		// OP_0 | Sig0 | Sig1 ... | RedeemScript
		// OP_0 for Multisig off-by-one error.
		if len(signs) > 0 {
			builder.AddOp(xvm.OP_0)
		}
		for _, sign := range signs {
			builder.AddData(sign.Signature)
		}
		return builder.AddData(redeemScript).Script()
	default:
		return nil, xerror.NewError(Errors, ER_SCRIPT_TYPE_UNKNOWN, class)
	}
}
