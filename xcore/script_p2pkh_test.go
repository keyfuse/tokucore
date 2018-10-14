// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokublock/tokucore/xvm"
)

func TestScriptP2PKH(t *testing.T) {
	outputScriptBytes := []byte{0x76, 0xa9, 0x14, 0x14, 0x83, 0x6d, 0xbe, 0x7f, 0x38, 0xc5, 0xac, 0x3d, 0x49, 0xe8, 0xd7, 0x90, 0xaf, 0x80, 0x8a, 0x4e, 0xe9, 0xed, 0xcf, 0x88, 0xac}
	outputScriptString := "OP_DUP OP_HASH160 OP_DATA_20 14836dbe7f38c5ac3d49e8d790af808a4ee9edcf OP_EQUALVERIFY OP_CHECKSIG"

	hex, _ := hex.DecodeString("14836dbe7f38c5ac3d49e8d790af808a4ee9edcf")
	script := NewPayToPubKeyHashScript(hex)
	locking, err := script.GetLockingScriptBytes()
	assert.Nil(t, err)
	assert.Equal(t, outputScriptBytes, locking)
	assert.Equal(t, outputScriptString, xvm.DisasmString(locking))

	assert.NotNil(t, script.GetAddress())
}
