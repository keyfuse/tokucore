// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xbase"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xcrypto/ripemd160"
	"github.com/tokublock/tokucore/xerror"
)

// Address -- an interface type for any type of destination a transaction.
// output may spent to, includes:
// 1. pay-to-pubkey (P2PK)
// 2. pay-to-pubkey-hash (P2PKH)
// 3. pay-to-script-hash (P2SH)
type Address interface {
	// ToString returns the string of the address with base58 encoding.
	ToString(net *network.Network) string

	Hash160() []byte
}

// DecodeAddress -- decode the string address and returns the Address with a known address type.
func DecodeAddress(addr string, net *network.Network) (Address, error) {
	decoded, netID, err := xbase.Base58CheckDecode(addr)
	if err != nil {
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
func NewPayToPubKeyHashAddress(hash160 []byte) Address {
	return &PayToPubKeyHashAddress{
		pubKeyHash: hash160,
	}
}

// ToString -- the implementation method for xcore.Address interface.
func (a *PayToPubKeyHashAddress) ToString(net *network.Network) string {
	return xbase.Base58CheckEncode(a.pubKeyHash[:ripemd160.Size], net.PubKeyHashAddrID)
}

// Hash160 -- the address hash160 bytes.
func (a *PayToPubKeyHashAddress) Hash160() []byte {
	return a.pubKeyHash[:ripemd160.Size]
}

// *******************************************
// PayToScriptHashAddress(P2SH)
// *******************************************

// PayToScriptHashAddress --
type PayToScriptHashAddress struct {
	scriptHash []byte
}

// NewPayToScriptHashAddress -- creates a new PayToScriptHashAddress.
func NewPayToScriptHashAddress(hash160 []byte) Address {
	return &PayToScriptHashAddress{
		scriptHash: hash160,
	}
}

// ToString -- the implementation method for xcore.Address interface.
func (a *PayToScriptHashAddress) ToString(net *network.Network) string {
	return xbase.Base58CheckEncode(a.scriptHash[:ripemd160.Size], net.ScriptHashAddrID)
}

// Hash160 -- returns the address hash160 bytes>
func (a *PayToScriptHashAddress) Hash160() []byte {
	return a.scriptHash[:ripemd160.Size]
}
