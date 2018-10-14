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
	"github.com/tokublock/tokucore/xcrypto"
)

func TestAddressP2PKH(t *testing.T) {
	net := network.MainNet
	{
		hexstr := "f6889b21b5540353a29ed18c45ea0031280c42cf"
		addr := "1PUYsjwfNmX64wS368ZR5FMouTtUmvtmTY"
		hex, _ := hex.DecodeString(hexstr)
		address := NewPayToPubKeyHashAddress(hex)
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
		address := NewPayToPubKeyHashAddress(hex)
		assert.Nil(t, address)
	}
}

// P2PKH example
func TestAddressP2PKHExample(t *testing.T) {
	seed := []byte("this.is.alice.seed.at.2018")
	key := NewHDKey(seed)
	pubkey := key.PublicKey()
	{
		uncompressed := pubkey.SerializeUncompressed()
		hash := xcrypto.Hash160(uncompressed)
		t.Log("p2pkh.uncompressed")
		t.Logf("\tpubkey:%X, len:%v", uncompressed, len(uncompressed))
		t.Logf("\thash:%X", hash)
		addr := NewPayToPubKeyHashAddress(hash)
		t.Logf("\taddress.testnet:%v", addr.ToString(network.TestNet))
		t.Logf("\taddress.mainnet:%v", addr.ToString(network.MainNet))
	}

	{
		compressed := pubkey.Serialize()
		hash := xcrypto.Hash160(compressed)
		t.Log("p2pkh.compressed")
		t.Logf("\tpubkey:%X, len:%v", compressed, len(compressed))
		t.Logf("\thash:%X", hash)
		addr := NewPayToPubKeyHashAddress(hash)
		t.Logf("\taddress.testnet:%v", addr.ToString(network.TestNet))
		t.Logf("\taddress.mainnet:%v", addr.ToString(network.MainNet))
	}
}
