// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokublock/tokucore/xcrypto"
)

func readTests(testfile string) ([][]interface{}, error) {
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

func TestEngine(t *testing.T) {
	t.Parallel()

	tests, err := readTests("testdata/vm.json")
	assert.Nil(t, err)

	for i, test := range tests {
		tst, ok := test[0].([]interface{})
		if !ok {
			continue
		}

		tname := tst[0]
		t.Logf("test.case[#%v-%v]", i, tname)
		lockingstr, ok := tst[1].(string)
		assert.True(t, ok)
		locking, err := NewScriptBuilder().Load(lockingstr).Script()
		assert.Nil(t, err)
		unlockingstr, ok := tst[2].(string)
		assert.True(t, ok)
		unlocking, err := NewScriptBuilder().Load(unlockingstr).Script()
		assert.Nil(t, err)

		errstr, ok := tst[3].(string)
		assert.True(t, ok)

		debug, ok := tst[4].(string)
		assert.True(t, ok)

		engine := NewEngine()

		if debug == "true" {
			engine.EnableDebug()
		} else {
			engine.DisableDebug()
		}
		// Hash function.
		hasherFn := func(hashType byte) []byte {
			return xcrypto.DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})
		}
		engine.SetHasher(hasherFn)

		// Verifier function.
		verifierFn := func(hash []byte, signature []byte, pubkey []byte) error {
			pub, err := xcrypto.PubKeyFromBytes(pubkey)
			if err != nil {
				return err
			}
			return xcrypto.Verify(hash, signature, pub)
		}
		engine.SetVerifier(verifierFn)

		err = engine.Verify(unlocking, locking)
		engine.PrintTrace()
		if errstr != "" {
			assert.NotNil(t, err)
			assert.Equal(t, errstr, err.Error())
		} else {
			if err != nil {
				t.Fatalf("#%s:%v", tname, err)
			}
		}
	}
}

func TestEngineExecute(t *testing.T) {
	engine := NewEngine()
	engine.EnableDebug()
	script, err := NewScriptBuilder().AddOp(OP_1).AddOp(OP_2).AddOp(OP_ADD).AddOp(OP_3).AddOp(OP_EQUAL).Script()
	assert.Nil(t, err)
	err = engine.Execute(script)
	assert.Nil(t, err)
}
