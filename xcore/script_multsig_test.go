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
	k1 := xcrypto.PrvKeyFromBytes([]byte{0x01})
	k2 := xcrypto.PrvKeyFromBytes([]byte{0x02})
	k3 := xcrypto.PrvKeyFromBytes([]byte{0x03})
	keys = append(keys, k1, k2, k3)

	for _, k := range keys {
		pk := k.PubKey().Serialize()
		pubkeys = append(pubkeys, pk)

		_, err := xcrypto.Sign(hash, k)
		assert.Nil(t, err)
	}

	// 2-of-3
	redeem, err := NewPayToMultiSigScript(2, pubkeys...).GetLockingScriptBytes()
	assert.Nil(t, err)
	t.Logf("2-of-3.script:%s", xvm.DisasmString(redeem))

	p2sh, err := NewPayToScriptHashScript(xcrypto.Hash160(redeem)).GetLockingScriptBytes()
	assert.Nil(t, err)
	t.Logf("locking:%s", xvm.DisasmString(p2sh))
}
