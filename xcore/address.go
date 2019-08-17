// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcore

import (
	"strings"

	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xbase"
	"github.com/keyfuse/tokucore/xerror"
)

// Address -- an interface type for any type of destination a transaction.
// output may spent to, includes:
// 1. pay-to-pubkey-hash (P2PKH)
// 2. pay-to-script-hash (P2SH)
// 3. pay-to-witness-pubkey-hash (P2WPKH)
type Address interface {
	// ToString returns the string of the address with base58 encoding.
	ToString(net *network.Network) string

	// Hash160 -- hash160 of this address.
	Hash160() []byte

	// LockingScript -- the locking script of this address.
	LockingScript() ([]byte, error)
}

// DecodeAddress -- decode the string address and returns the Address with a known address type.
func DecodeAddress(addr string, net *network.Network) (Address, error) {
	// Bech32 encoded segwit addresses start with a human-readable part (hrp) followed by '1'(bc or tb)
	oneIndex := strings.IndexByte(addr, '1')
	if oneIndex == 2 {
		prefix := addr[:oneIndex]
		if prefix == net.Bech32HRPSegwit {
			_, version, witnessProgram, err := xbase.WitnessDecode(addr)
			if err != nil {
				return nil, err
			}
			switch version {
			case 0x00:
				switch len(witnessProgram) {
				case 20:
					return NewPayToWitnessV0PubKeyHashAddress(witnessProgram), nil
				case 32:
					return NewPayToWitnessV0ScriptHashAddress(witnessProgram), nil
				}
			default:
				return nil, xerror.NewError(Errors, ER_ADDRESS_WITNESS_VERSION_UNSUPPORTED, version)
			}
		}
	}

	decoded, netID, err := xbase.Base58CheckDecode(addr)
	if err != nil {
		return nil, xerror.NewError(Errors, ER_ADDRESS_FORMAT_MALFORMED, addr)
	}
	switch len(decoded) {
	case 20: // RIPEMD160 size, P2PKH or P2SH
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
