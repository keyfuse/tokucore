// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScriptReaderParser(t *testing.T) {
	hex, err := hex.DecodeString("14836dbe7f38c5ac3d49e8d790af808a4ee9edcf")
	assert.Nil(t, err)

	// Build.
	script, err := NewScriptBuilder().AddOp(OP_DUP).AddOp(OP_HASH160).AddData(hex).AddOp(OP_EQUALVERIFY).AddOp(OP_CHECKSIG).Script()
	assert.Nil(t, err)

	// Parsed.
	pops, err := NewScriptReader(script).AllInstructions()
	assert.Nil(t, err)
	for _, pop := range pops {
		t.Logf("--: %+v", pop.op.name)
	}
}

func TestScriptReaderRemoveOpCode(t *testing.T) {
	data0 := bytes.Repeat([]byte{0x01}, 0xff-1)
	data1 := bytes.Repeat([]byte{0x01}, 0xff+1)
	data2 := bytes.Repeat([]byte{0x01}, 0xffff+1)

	script, err := NewScriptBuilder().AddOp(OP_DUP).AddOp(OP_HASH160).AddInt64(128).AddData(data0).AddData(data1).AddData(data2).AddOp(OP_EQUALVERIFY).AddOp(OP_CHECKSIG).Script()
	assert.Nil(t, err)

	removed, err := RemoveOpcode(script, OP_PUSHDATA1)
	assert.Nil(t, err)

	removed, err = RemoveOpcode(removed, OP_PUSHDATA2)
	assert.Nil(t, err)

	removed, err = RemoveOpcode(removed, OP_PUSHDATA4)
	assert.Nil(t, err)

	got := DisasmString(removed)
	want := "OP_DUP OP_HASH160 OP_DATA_2 8000 OP_EQUALVERIFY OP_CHECKSIG"
	assert.Equal(t, want, got)
}
