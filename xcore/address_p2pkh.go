// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcore

import (
	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xbase"
)

// *******************************************
// PayToPubKeyHashAddress(P2PKH)
// *******************************************

// PayToPubKeyHashAddress --
type PayToPubKeyHashAddress struct {
	pubKeyHash [20]byte
}

// NewPayToPubKeyHashAddress -- creates a new PayToPubKeyHashAddress.
func NewPayToPubKeyHashAddress(hash160 []byte) Address {
	if len(hash160) != 20 {
		return nil
	}

	var hash [20]byte
	copy(hash[:], hash160)
	return &PayToPubKeyHashAddress{
		pubKeyHash: hash,
	}
}

// ToString -- the implementation method for xcore.Address interface.
func (a *PayToPubKeyHashAddress) ToString(net *network.Network) string {
	return xbase.Base58CheckEncode(a.pubKeyHash[:], net.PubKeyHashAddrID)
}

// Hash160 -- the address hash160 bytes.
func (a *PayToPubKeyHashAddress) Hash160() []byte {
	return a.pubKeyHash[:]
}

// LockingScript -- the address locking script.
func (a *PayToPubKeyHashAddress) LockingScript() ([]byte, error) {
	return NewPayToPubKeyHashScript(a.Hash160()).GetRawLockingScriptBytes()
}
