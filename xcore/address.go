// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"github.com/tokublock/tokucore/xbase/base58"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xerror"
)

// Address -- an interface type for any type of destination a transaction.
// output may spent to, includes:
// 1. pay-to-pubkey (P2PK)
// 2. pay-to-pubkey-hash (P2PKH)
// 3. pay-to-script-hash (P2SH)
type Address interface {
	// ToString returns the string of the address with base58 encoding.
	ToString(net *NetworkParams) string

	Hash160() []byte
}

// DecodeAddress -- decode the string address and returns the Address with a known address type.
func DecodeAddress(addr string, net *NetworkParams) (Address, error) {
	decoded, netID, err := base58.CheckDecode(addr)
	if err != nil {
		if err == base58.ErrChecksum {
			return nil, xerror.NewError(Errors, ER_ADDRESS_CHECKSUM_MISMATCH)
		}
		return nil, xerror.NewError(Errors, ER_ADDRESS_FORMAT_MALFORMED, addr)
	}
	switch len(decoded) {
	case xcrypto.Ripemd160Size(): // RIPEMD160 size, P2PKH or P2SH
		switch netID {
		case net.PubKeyHashAddrID:
			return NewPayToPubKeyHashAddress(decoded), nil
		case net.ScriptHashAddrID:
			return NewPayToScriptHashAddress(decoded), nil
		default:
			return nil, xerror.NewError(Errors, ER_ADDRESS_TYPE_UNKNOWN, netID)
		}
	default:
		return nil, xerror.NewError(Errors, ER_ADDRESS_SIZE_MALFORMED, len(decoded))
	}
}

// *******************************************
// PayToPubKeyHashAddress(P2PKH)
// *******************************************

// PayToPubKeyHashAddress --
type PayToPubKeyHashAddress struct {
	pubKeyHash []byte
}

// NewPayToPubKeyHashAddress -- creates a new PayToPubKeyHashAddress.
func NewPayToPubKeyHashAddress(pubKeyHash []byte) Address {
	return &PayToPubKeyHashAddress{
		pubKeyHash: pubKeyHash,
	}
}

// ToString -- the implementation method for xcore.Address interface.
func (a *PayToPubKeyHashAddress) ToString(net *NetworkParams) string {
	return base58.CheckEncode(a.pubKeyHash, net.PubKeyHashAddrID)
}

// Hash160 -- the address hash160 bytes.
func (a *PayToPubKeyHashAddress) Hash160() []byte {
	return a.pubKeyHash
}

// *******************************************
// PayToScriptHashAddress(P2SH)
// *******************************************

// PayToScriptHashAddress --
type PayToScriptHashAddress struct {
	scriptHash []byte
}

// NewPayToScriptHashAddress -- creates a new PayToScriptHashAddress.
func NewPayToScriptHashAddress(scriptHash []byte) Address {
	return &PayToScriptHashAddress{
		scriptHash: scriptHash,
	}
}

// ToString -- the implementation method for xcore.Address interface.
func (a *PayToScriptHashAddress) ToString(net *NetworkParams) string {
	return base58.CheckEncode(a.scriptHash, net.ScriptHashAddrID)
}

// Hash160 -- returns the address hash160 bytes>
func (a *PayToScriptHashAddress) Hash160() []byte {
	return a.scriptHash
}
