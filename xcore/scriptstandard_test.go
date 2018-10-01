// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xerror"
	"github.com/tokublock/tokucore/xvm"
)

func TestScriptDisasm(t *testing.T) {
	scripts := []struct {
		name string
		hex  string
	}{
		{
			name: "#1",
			hex:  "522103c8727ce35f1c93eb0be21406ee9a923c89219fe9c9e8504c8314a6a22d1295c02103c74dc710c407d7db6e041ee212d985cd2826d93f806ed44912b9a1da691c977352ae",
		},
		{
			name: "#2",
			hex:  "0020a8f44467bf171d51499153e01c0bd6291109fc38bd21b3c3224c9dc6b57590df",
		},
		{
			name: "#3",
			hex:  "483045022100ea6646844e0a228eaef712af567ba2a377e5505a66c07c74c195c01c50c5cff3022071f83ff8edc63b1f3b5c90e89ae16e20ad0ef7b720e0e231c7d1e6ebbed83da20121034f754f4d4180716540935b81e74a46712ce785cec54e73b278d7c2d6d90f9895",
		},
	}

	for _, test := range scripts {
		_, err := hex.DecodeString(test.hex)
		assert.Nil(t, err)
	}
}

func TestScript(t *testing.T) {
	tests := []struct {
		name               string
		fn                 func([]byte) Script
		outputScriptBytes  []byte
		outputScriptString string
	}{
		{
			name:               "PayToPubKeyHashScript",
			fn:                 NewPayToPubKeyHashScript,
			outputScriptBytes:  []byte{0x76, 0xa9, 0x14, 0x14, 0x83, 0x6d, 0xbe, 0x7f, 0x38, 0xc5, 0xac, 0x3d, 0x49, 0xe8, 0xd7, 0x90, 0xaf, 0x80, 0x8a, 0x4e, 0xe9, 0xed, 0xcf, 0x88, 0xac},
			outputScriptString: "OP_DUP OP_HASH160 OP_DATA_20 14836dbe7f38c5ac3d49e8d790af808a4ee9edcf OP_EQUALVERIFY OP_CHECKSIG",
		},
		{
			name:               "PayToScriptHashScript",
			fn:                 NewPayToScriptHashScript,
			outputScriptBytes:  []byte{0xa9, 0x14, 0x14, 0x83, 0x6d, 0xbe, 0x7f, 0x38, 0xc5, 0xac, 0x3d, 0x49, 0xe8, 0xd7, 0x90, 0xaf, 0x80, 0x8a, 0x4e, 0xe9, 0xed, 0xcf, 0x87},
			outputScriptString: "OP_HASH160 OP_DATA_20 14836dbe7f38c5ac3d49e8d790af808a4ee9edcf OP_EQUAL",
		},
	}

	hex, _ := hex.DecodeString("14836dbe7f38c5ac3d49e8d790af808a4ee9edcf")
	for _, test := range tests {
		t.Logf("test:%v", test.name)
		script := test.fn(hex)
		locking, err := script.GetLockingScriptBytes()
		assert.Nil(t, err)
		assert.Equal(t, test.outputScriptBytes, locking)
		assert.Equal(t, test.outputScriptString, xvm.DisasmString(locking))
	}
}

// TestMultiSigScript ensures the MultiSigScript function returns the expected
// scripts and errors.
func TestMultiSigScript(t *testing.T) {
	t.Parallel()

	pkComressed1, _ := hex.DecodeString("02192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b4")
	pkComressed2, _ := hex.DecodeString("03b0bd634234abbb1ba1e986e884185c61cf43e001f9137f23c2c409273eb16e65")
	pkUncompressed, _ := hex.DecodeString("0411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3")

	tests := []struct {
		keys      [][]byte
		nrequired int
		expected  string
		err       error
	}{
		{
			[][]byte{
				pkComressed1,
				pkComressed2,
			},
			1,
			"OP_1 OP_DATA_33 02192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b4 OP_DATA_33 03b0bd634234abbb1ba1e986e884185c61cf43e001f9137f23c2c409273eb16e65 OP_2 OP_CHECKMULTISIG",
			nil,
		},
		{
			[][]byte{
				pkComressed1,
				pkComressed2,
			},
			2,
			"OP_2 OP_DATA_33 02192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b4 OP_DATA_33 03b0bd634234abbb1ba1e986e884185c61cf43e001f9137f23c2c409273eb16e65 OP_2 OP_CHECKMULTISIG",
			nil,
		},
		{
			[][]byte{
				pkComressed1,
				pkComressed2,
			},

			3,
			"",
			xerror.NewError(Errors, ER_SCRIPT_STANDARD_PUBKEYS_LE_NREQUIRED, 2, 3),
		},
		{
			[][]byte{
				pkUncompressed,
			},
			1,
			"OP_1 OP_DATA_65 0411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3 OP_1 OP_CHECKMULTISIG",
			nil,
		},
		{
			[][]byte{
				pkComressed2,
			},
			2,
			"",
			xerror.NewError(Errors, ER_SCRIPT_STANDARD_PUBKEYS_LE_NREQUIRED, 1, 2),
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		script := NewPayToMultiSigScript(test.nrequired, test.keys...)
		locking, err := script.GetLockingScriptBytes()
		if err != nil {
			assert.Equal(t, test.err.Error(), err.Error())
		}
		asm := xvm.DisasmString(locking)
		if asm != test.expected {
			t.Errorf("MultiSigScript #%d\n\tgot:  %s\n\twant: %s", i, asm, test.expected)
			continue
		}
	}
}

func TestMultiSigScriptDemo(t *testing.T) {
	hash := xcrypto.DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	var keys []*xcrypto.PrivateKey
	var pubkeys [][]byte
	var pubsigs []PubKeySign
	k1 := xcrypto.PrvKeyFromBytes([]byte{0x01})
	k2 := xcrypto.PrvKeyFromBytes([]byte{0x02})
	k3 := xcrypto.PrvKeyFromBytes([]byte{0x03})
	keys = append(keys, k1, k2, k3)

	for _, k := range keys {
		pk := k.PubKey().Serialize()
		pubkeys = append(pubkeys, pk)

		sig, err := xcrypto.Sign(hash, k)
		assert.Nil(t, err)
		pubsigs = append(pubsigs, PubKeySign{
			pk,
			sig,
		})
	}

	// 2-of-3
	redeem, err := NewPayToMultiSigScript(2, pubkeys...).GetLockingScriptBytes()
	assert.Nil(t, err)
	t.Logf("2-of-3.script:%s", xvm.DisasmString(redeem))

	p2sh, err := NewPayToScriptHashScript(xcrypto.Hash160(redeem)).GetLockingScriptBytes()
	assert.Nil(t, err)
	t.Logf("locking:%s", xvm.DisasmString(p2sh))

	// Unlocking 01.
	{
		pubsigs01 := []PubKeySign{pubsigs[0], pubsigs[1]}
		unlocking, err := BuildUnlockingScriptBytes(p2sh, redeem, pubsigs01)
		assert.Nil(t, err)
		t.Logf("unlocking01:%s", xvm.DisasmString(unlocking))
	}

	// Unlocking 02.
	{
		pubsigs02 := []PubKeySign{pubsigs[0], pubsigs[2]}
		unlocking, err := BuildUnlockingScriptBytes(p2sh, redeem, pubsigs02)
		assert.Nil(t, err)
		t.Logf("unlocking02:%s", xvm.DisasmString(unlocking))
	}

	// Unlocking 10.
	{
		pubsigs10 := []PubKeySign{pubsigs[1], pubsigs[0]}
		unlocking, err := BuildUnlockingScriptBytes(p2sh, redeem, pubsigs10)
		assert.Nil(t, err)
		t.Logf("unlocking10:%s", xvm.DisasmString(unlocking))
	}
}
