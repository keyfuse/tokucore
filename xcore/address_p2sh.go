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
// PayToScriptHashAddress(P2SH)
// *******************************************

// PayToScriptHashAddress --
type PayToScriptHashAddress struct {
	scriptHash [20]byte
}

// NewPayToScriptHashAddress -- creates a new PayToScriptHashAddress.
func NewPayToScriptHashAddress(hash160 []byte) Address {
	if len(hash160) != 20 {
		return nil
	}

	var hash [20]byte
	copy(hash[:], hash160)
	return &PayToScriptHashAddress{
		scriptHash: hash,
	}
}

// ToString -- the implementation method for xcore.Address interface.
func (a *PayToScriptHashAddress) ToString(net *network.Network) string {
	return xbase.Base58CheckEncode(a.scriptHash[:], net.ScriptHashAddrID)
}

// Hash160 -- returns the address hash160 bytes>
func (a *PayToScriptHashAddress) Hash160() []byte {
	return a.scriptHash[:]
}

// LockingScript -- the address locking script.
func (a *PayToScriptHashAddress) LockingScript() ([]byte, error) {
	return NewPayToScriptHashScript(a.Hash160()).GetRawLockingScriptBytes()
}
