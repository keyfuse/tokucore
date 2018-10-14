// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"strings"

	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xbase"
	"github.com/tokublock/tokucore/xerror"
)

// Address -- an interface type for any type of destination a transaction.
// output may spent to, includes:
// 1. pay-to-pubkey-hash (P2PKH)
// 2. pay-to-script-hash (P2SH)
// 3. pay-to-witness-pubkey-hash (P2WPKH)
type Address interface {
	// ToString returns the string of the address with base58 encoding.
	ToString(net *network.Network) string

	Hash160() []byte
}

// DecodeAddress -- decode the string address and returns the Address with a known address type.
func DecodeAddress(addr string, net *network.Network) (Address, error) {
	// Bech32 encoded segwit addresses start with a human-readable part (hrp) followed by '1'(bc or tb)
	oneIndex := strings.IndexByte(addr, '1')
	if oneIndex == 2 {
		prefix := addr[:oneIndex]
		if prefix == net.Bech32HRPSegwit {
			_, version, witnessProgram, err := WitnessAddressDecode(addr)
			if err != nil {
				return nil, err
			}
			if version != 0 {
				return nil, xerror.NewError(Errors, ER_ADDRESS_WITNESS_VERSION_UNSUPPORTED, version)
			}
			switch len(witnessProgram) {
			case 20:
				return NewPayToWitnessPubKeyHashAddress(witnessProgram), nil
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
