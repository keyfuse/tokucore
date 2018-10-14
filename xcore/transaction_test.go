// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokublock/tokucore/xbase"
)

// https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki
func TestTransactionRaw(t *testing.T) {
	outputs := []struct {
		txid   string
		index  uint32
		script string
		value  uint64
	}{
		{
			txid:   "0a4381c05136c0cb44886a5df7c26f1930bcc2c12e00ec60e027c4378d7d8c2e",
			index:  1,
			script: "a914203736c3c06053896d7041ce8f5bae3df76cc49187",
			value:  0.5 * 1e8,
		},
		{
			txid:   "2c4df245d00b491bdf24965adbbccdaa7f62ccac933d3e9377f336c60c4ea096",
			index:  0,
			script: "a914f3ba8a120d960ae07d1dbe6f0c37fb4c926d76d587",
			value:  2.0 * 1e8,
		},
	}
	tx := NewTransaction()
	tx.SetVersion(2)

	// Input.
	for _, out := range outputs {
		txhash, err := xbase.NewIDFromString(out.txid)
		assert.Nil(t, err)
		tx.AddInput(&TxIn{
			Hash:     txhash,
			Index:    out.index,
			Sequence: 0xffffffff,
		})
	}

	// Output.
	scriptHash, err := hex.DecodeString("b53bb0dc1db8c8d803e3e39f784d42e4737ffa0d")
	assert.Nil(t, err)
	lockingScript, err := NewPayToScriptHashScript(scriptHash).GetLockingScriptBytes()
	assert.Nil(t, err)
	tx.AddOutput(&TxOut{
		Value:  249900000,
		Script: lockingScript,
	})

	tx.SetLockTime(7)
	tx.SetSigHashType(SigHashAll)
}

func TestTransactions(t *testing.T) {
	tests, err := readTxTests("testdata/tx.json")
	assert.Nil(t, err)

	for _, test := range tests {
		// Inputs.
		inputs, ok := test[0].([]interface{})
		if !ok {
			continue
		}

		txs, ok := test[1].([]interface{})
		if !ok {
			continue
		}

		// Witness.
		witness, ok := test[2].(bool)
		if !ok {
			continue
		}

		// Name.
		tName, ok := test[3].(string)
		if !ok {
			continue
		}
		t.Logf("%s", tName)

		var txid string
		var serializedTx []byte
		for _, tx := range txs {
			data, ok := tx.([]interface{})
			if !ok {
				t.Errorf("bad.tx.data")
				continue
			}

			// Txid.
			txid, ok = data[0].(string)
			if !ok {
				t.Errorf("bad.txid.hex")
				continue
			}

			// Serialized tx hex.
			serializedTxHex, ok := data[1].(string)
			if !ok {
				continue
			}
			serializedTx, err = hex.DecodeString(serializedTxHex)
			assert.Nil(t, err)
		}

		tx := NewTransaction()
		if witness {
			err = tx.Deserialize(serializedTx)
		} else {
			err = tx.DeserializeNoWitness(serializedTx)
		}
		assert.Nil(t, err)

		for i, iinput := range inputs {
			input, ok := iinput.([]interface{})
			if !ok {
				t.Errorf("bad.input.hex")
				continue
			}

			// Amount.
			amount, ok := input[1].(float64)
			if !ok {
				t.Errorf("bad.amount")
				continue
			}

			// Locking.
			lockingHex, ok := input[2].(string)
			if !ok {
				t.Errorf("bad.locking.hex")
				continue
			}
			locking, err := hex.DecodeString(lockingHex)
			assert.Nil(t, err)

			// Redeem hex.
			redeemHex, ok := input[3].(string)
			if !ok {
				t.Errorf("bad.redeem.hex")
				continue
			}
			var redeem []byte
			if redeemHex != "" {
				redeem, err = hex.DecodeString(redeemHex)
				assert.Nil(t, err)
			}
			tx.SetTxIn(i, uint64(amount), locking, redeem)
		}

		// Debug.
		t.Logf("%v", tx.ToString())

		// Verify.
		if err = tx.Verify(); err != nil {
			t.Fatalf("%s.verify.failed.err:%v", tName, err)
		}

		// Txid check.
		txid1 := tx.ID()
		assert.Equal(t, txid, txid1)
	}
}
