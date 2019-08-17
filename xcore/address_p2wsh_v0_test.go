// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcore

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xvm"
)

func TestAddressP2WSH(t *testing.T) {
	net := network.MainNet
	{
		hexstr := "b41568b8c8af2fd75f370574b629de457ecbb7bb377f9f1a602d8dfb42bc8962"
		addr := "bc1qks2k3wxg4uhawhehq46tv2w7g4lvhdamxale7xnq9kxlks4u393qmltg4f"
		hex, _ := hex.DecodeString(hexstr)
		address := NewPayToWitnessV0ScriptHashAddress(hex)
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
		address := NewPayToWitnessV0ScriptHashAddress(hex)
		assert.Nil(t, address)
	}
}

// P2WSH example.
func TestAddressP2WSHExample1(t *testing.T) {
	pubkeys := MockPublicKeys(3)
	multisig := NewPayToMultiSigScript(2, pubkeys[0], pubkeys[1], pubkeys[2])
	witnessScript, err := multisig.GetLockingScriptBytes()
	assert.Nil(t, err)
	t.Logf("p2wsh.witness.script:%v", xvm.DisasmString(witnessScript))
	witnessScriptHash := sha256.Sum256(witnessScript)
	addr := NewPayToWitnessV0ScriptHashAddress(witnessScriptHash[:])
	t.Logf("p2wsh.address.testnte:%v", addr.ToString(network.TestNet))
	t.Logf("p2wsh.address.mainnet:%v", addr.ToString(network.MainNet))
}
