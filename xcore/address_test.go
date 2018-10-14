// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xbase"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xerror"
)

func TestAddressDecode(t *testing.T) {
	tests := []struct {
		net  *network.Network
		addr string
		err  error
	}{
		{
			net:  network.MainNet,
			addr: "1PUYsjwfNmX64wS368ZR5FMouTtUmvtmTY",
			err:  nil,
		},
		{
			net:  network.TestNet,
			addr: "n3zWAo2eBnxLr3ueohXnuAa8mTVBhxmPhq",
			err:  nil,
		},

		{
			net:  network.MainNet,
			addr: "3QAZoHS6vfqUA78UDEE1Vsik3zBCLfybyE",
			err:  nil,
		},
		{
			net:  network.TestNet,
			addr: "2NFims2N8Y8LpMtm1tMqt7pi1GLPN8pBgBc",
			err:  nil,
		},
		{
			net:  network.TestNet,
			addr: "2NFims2N8Y8LpMtm1tMqt7pi1GLPN8pBgBx",
			err:  xerror.NewError(Errors, ER_ADDRESS_FORMAT_MALFORMED, "2NFims2N8Y8LpMtm1tMqt7pi1GLPN8pBgBx"),
		},
	}

	for _, test := range tests {
		// Get the address interface.
		address, err := DecodeAddress(test.addr, test.net)
		if test.err == nil {
			assert.Nil(t, err)

			// Get the address string.
			addr := address.ToString(test.net)
			assert.Equal(t, test.addr, addr)
		} else {
			t.Logf("%+v", err)
			assert.Equal(t, test.err.Error(), err.Error())
		}
	}
}

func TestAddressExdous(t *testing.T) {
	for i := 0; i < 5; i++ {
		key, err := NewHDKeyRand()
		assert.Nil(t, err)
		t.Logf("priv:%v", key.ToString(network.TestNet))
		addr := NewPayToPubKeyHashAddress(key.PublicKey().Hash160())
		t.Logf("addr:%v", addr.ToString(network.TestNet))
	}
}

// Vanity private key example.
func TestAddressVanity(t *testing.T) {
	base58Prv := "5JLUmjZiirgziDmWmNprPsNx8DYwfecUNk1FQXmDPaoKB36fX1o"
	decoded, _, err := xbase.Base58CheckDecode(base58Prv)
	assert.Nil(t, err)
	t.Logf("%X, len:%v", decoded, len(decoded))

	pub := xcrypto.PrvKeyFromBytes(decoded).PubKey()
	hash := xcrypto.Hash160(pub.SerializeUncompressed())
	t.Logf("hash:%X", hash)
	got := NewPayToPubKeyHashAddress(hash)
	want := "1LoveRg5t2NCDLUZh6Q8ixv74M5YGVxXaN"
	assert.Equal(t, want, got.ToString(network.MainNet))
}
