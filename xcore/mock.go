// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"

	"github.com/tokublock/tokucore/xcrypto"
)

// MockP2PKHCoin -- mock p2pkh coin for tests.
func MockP2PKHCoin(hdKey *HDKey) *Coin {
	txid := make([]byte, 32)
	rand.Read(txid)

	script, _ := NewPayToPubKeyHashScript(hdKey.PublicKey().Hash160()).GetLockingScriptBytes()
	return &Coin{
		txID:   fmt.Sprintf("%x", txid),
		n:      0,
		value:  10000,
		script: fmt.Sprintf("%x", script),
	}
}

// MockP2SHCoin -- mocks p2sh coin for tests.
func MockP2SHCoin(alice *HDKey, bob *HDKey, redeem []byte) *Coin {
	//bobAlice := NewPayToMultiSigScript(2, alice.PublicKey().Serialize(), bob.PublicKey().Serialize())
	//bobAliceScript, _ := bobAlice.GetLockingScriptBytes()
	redeemScript := xcrypto.Hash160(redeem)

	txid := make([]byte, 32)
	rand.Read(txid)
	script, _ := NewPayToScriptHashScript(redeemScript).GetLockingScriptBytes()
	return &Coin{
		txID:   fmt.Sprintf("%x", txid),
		n:      1,
		value:  20000,
		script: fmt.Sprintf("%x", script),
	}
}

// MockPublicKeys -- mock publickey(serialize compressed) for tests.
func MockPublicKeys(num int) [][]byte {
	var keys [][]byte
	for i := 0; i < num; i++ {
		hdkey, _ := NewHDKeyRand()
		key := hdkey.PublicKey().Serialize()
		keys = append(keys, key)
	}
	return keys
}

// readTxTests -- reads datas from testdata/tx.json
func readTxTests(testfile string) ([][]interface{}, error) {
	file, err := ioutil.ReadFile(testfile)
	if err != nil {
		return nil, err
	}

	var tests [][]interface{}
	if err := json.Unmarshal(file, &tests); err != nil {
		return nil, err
	}
	return tests, nil
}
