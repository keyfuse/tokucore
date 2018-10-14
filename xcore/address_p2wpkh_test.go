// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokublock/tokucore/network"
)

func TestAddressP2WPKH(t *testing.T) {
	net := network.MainNet
	{
		hexstr := "f6889b21b5540353a29ed18c45ea0031280c42cf"
		addr := "bc1q76yfkgd42sp48g576xxyt6sqxy5qcsk0t7f50w"
		hex, _ := hex.DecodeString(hexstr)
		address := NewPayToWitnessPubKeyHashAddress(hex)
		assert.Equal(t, addr, address.ToString(net))
		address.Hash160()

		decode, err := DecodeAddress(addr, net)
		assert.Nil(t, err)
		assert.Equal(t, address, decode)
	}

	// nil.
	{
		hexstr := "f6889b21b5540353a29ed18c45ea0031280c42"
		hex, _ := hex.DecodeString(hexstr)
		address := NewPayToWitnessPubKeyHashAddress(hex)
		assert.Nil(t, address)
	}
}

// P2WPKH example.
func TestAddressP2WPKHExample(t *testing.T) {
	seed := []byte("this.is.alice.seed.at.2018")
	key := NewHDKey(seed)
	pubkey := key.PublicKey()
	hash := pubkey.Hash160()
	t.Logf("\thash:%X", hash)
	addr := NewPayToWitnessPubKeyHashAddress(hash)
	t.Logf("\taddress.testnet:%v", addr.ToString(network.TestNet))
	t.Logf("\taddress.mainnet:%v", addr.ToString(network.MainNet))
}
