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
	"github.com/tokublock/tokucore/xbase"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xerror"
	"github.com/tokublock/tokucore/xvm"
)

func TestAddress(t *testing.T) {
	tests := []struct {
		name   string
		hex    string
		fn     func([]byte) Address
		net    *network.Network
		addr   string
		script string
		id     string
	}{
		{
			name:   "PayToPubKeyHashAddress.MainNet",
			hex:    "f6889b21b5540353a29ed18c45ea0031280c42cf",
			fn:     NewPayToPubKeyHashAddress,
			net:    network.MainNet,
			addr:   "1PUYsjwfNmX64wS368ZR5FMouTtUmvtmTY",
			script: "OP_DUP OP_HASH160 OP_DATA_20 f6889b21b5540353a29ed18c45ea0031280c42cf OP_EQUALVERIFY OP_CHECKSIG",
			id:     "f6889b21b5540353a29ed18c45ea0031280c42cf",
		},
		{
			name:   "PayToPubKeyHashAddress.TestNet",
			hex:    "f6889b21b5540353a29ed18c45ea0031280c42cf",
			fn:     NewPayToPubKeyHashAddress,
			net:    network.TestNet,
			addr:   "n3zWAo2eBnxLr3ueohXnuAa8mTVBhxmPhq",
			script: "OP_DUP OP_HASH160 OP_DATA_20 f6889b21b5540353a29ed18c45ea0031280c42cf OP_EQUALVERIFY OP_CHECKSIG",
			id:     "f6889b21b5540353a29ed18c45ea0031280c42cf",
		},

		{
			name:   "PayToScriptHashAddress.MainNet",
			hex:    "f6889b21b5540353a29ed18c45ea0031280c42cf",
			fn:     NewPayToScriptHashAddress,
			net:    network.MainNet,
			addr:   "3QAZoHS6vfqUA78UDEE1Vsik3zBCLfybyE",
			script: "OP_HASH160 OP_DATA_20 f6889b21b5540353a29ed18c45ea0031280c42cf OP_EQUAL",
			id:     "f6889b21b5540353a29ed18c45ea0031280c42cf",
		},
		{
			name:   "PayToScriptHashAddress.TestNet",
			hex:    "f6889b21b5540353a29ed18c45ea0031280c42cf",
			fn:     NewPayToScriptHashAddress,
			net:    network.TestNet,
			addr:   "2NFims2N8Y8LpMtm1tMqt7pi1GLPN8pBgBc",
			script: "OP_HASH160 OP_DATA_20 f6889b21b5540353a29ed18c45ea0031280c42cf OP_EQUAL",
			id:     "f6889b21b5540353a29ed18c45ea0031280c42cf",
		},
	}

	for _, test := range tests {
		t.Logf("test:%v", test.name)
		hex, _ := hex.DecodeString(test.hex)
		address := test.fn(hex)
		assert.Equal(t, test.addr, address.ToString(test.net))
	}
}

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
			assert.Equal(t, test.err, err)
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

// MultiSig example
func TestAddressP2SHExample1(t *testing.T) {
	pubkeys := MockPublicKeys(3)
	multisig := NewPayToMultiSigScript(2, pubkeys[0], pubkeys[1], pubkeys[2])
	redeemScript, err := multisig.GetLockingScriptBytes()
	assert.Nil(t, err)
	t.Logf("p2sh.redeem.script:%v", xvm.DisasmString(redeemScript))
	addr := multisig.GetAddress()
	t.Logf("p2sh.address.testnte:%v", addr.ToString(network.TestNet))
	t.Logf("p2sh.address.mainnet:%v", addr.ToString(network.MainNet))
}

// Unstandard Script example
func TestAddressP2SHExample2(t *testing.T) {
	// x + y = 7
	redeemScript, err := xvm.NewScriptBuilder().AddOp(xvm.OP_ADD).AddOp(xvm.OP_7).AddOp(xvm.OP_EQUAL).Script()
	assert.Nil(t, err)
	t.Logf("p2sh.redeem.script:%v", xvm.DisasmString(redeemScript))
	hash := xcrypto.Hash160(redeemScript)
	t.Logf("p2sh.redeem.script.hash:%X", hash)
	addr := NewPayToScriptHashAddress(hash)
	t.Logf("p2sh.address.testnte:%v", addr.ToString(network.TestNet))
	t.Logf("p2sh.address.mainnet:%v", addr.ToString(network.MainNet))
}

// Pay to:
// z + y = 9, z + x = 7, x + y = 8
//
// OP_3 OP_5 OP_4
// OP_3DUP
// OP_ADD OP_9 OP_EQUALVERIFY
// OP_ADD OP_7 OP_EQUALVERIFY
// OP_ADD OP_8 OP_EQUALVERIFY
// OP_1
func TestAddressP2SHExample3(t *testing.T) {
	redeemScript, err := xvm.NewScriptBuilder().AddOp(xvm.OP_3DUP).
		AddOp(xvm.OP_ADD).AddOp(xvm.OP_9).AddOp(xvm.OP_EQUALVERIFY).
		AddOp(xvm.OP_ADD).AddOp(xvm.OP_7).AddOp(xvm.OP_EQUALVERIFY).
		AddOp(xvm.OP_ADD).AddOp(xvm.OP_8).AddOp(xvm.OP_EQUALVERIFY).
		AddOp(xvm.OP_1).
		Script()
	assert.Nil(t, err)
	t.Logf("p2sh.redeem.script:%v", xvm.DisasmString(redeemScript))
	hash := xcrypto.Hash160(redeemScript)
	t.Logf("p2sh.redeem.script.hash:%X", hash)
	addr := NewPayToScriptHashAddress(hash)
	t.Logf("p2sh.address.testnte:%v", addr.ToString(network.TestNet))
	t.Logf("p2sh.address.mainnet:%v", addr.ToString(network.MainNet))
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
