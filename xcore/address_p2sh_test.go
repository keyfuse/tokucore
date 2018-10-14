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
	"github.com/tokublock/tokucore/xvm"
)

func TestAddressP2SH(t *testing.T) {
	net := network.MainNet
	{
		hexstr := "f6889b21b5540353a29ed18c45ea0031280c42cf"
		addr := "3QAZoHS6vfqUA78UDEE1Vsik3zBCLfybyE"
		hex, _ := hex.DecodeString(hexstr)
		address := NewPayToScriptHashAddress(hex)
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
		address := NewPayToScriptHashAddress(hex)
		assert.Nil(t, address)
	}
}

// MultiSig P2SH example
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
