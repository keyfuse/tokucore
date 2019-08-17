// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcore

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/keyfuse/tokucore/xvm"
)

func TestScriptP2WSH(t *testing.T) {
	outputScriptBytes := []byte{0x0, 0x20, 0xb4, 0x15, 0x68, 0xb8, 0xc8, 0xaf, 0x2f, 0xd7, 0x5f, 0x37, 0x5, 0x74, 0xb6, 0x29, 0xde, 0x45, 0x7e, 0xcb, 0xb7, 0xbb, 0x37, 0x7f, 0x9f, 0x1a, 0x60, 0x2d, 0x8d, 0xfb, 0x42, 0xbc, 0x89, 0x62}
	outputScriptString := "OP_0 OP_DATA_32 b41568b8c8af2fd75f370574b629de457ecbb7bb377f9f1a602d8dfb42bc8962"

	hex, _ := hex.DecodeString("b41568b8c8af2fd75f370574b629de457ecbb7bb377f9f1a602d8dfb42bc8962")
	script := NewPayToWitnessV0ScriptHashScript(hex)
	locking, err := script.GetRawLockingScriptBytes()
	assert.Nil(t, err)
	assert.Equal(t, outputScriptBytes, locking)
	assert.Equal(t, outputScriptString, xvm.DisasmString(locking))

	assert.NotNil(t, script.GetAddress())
}
