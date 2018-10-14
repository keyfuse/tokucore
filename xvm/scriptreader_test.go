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
	hex, err := hex.DecodeString("02483045022100df7b7e5cda14ddf91290e02ea10786e03eb11ee36ec02dd862fe9a326bbcb7fd02203f5b4496b667e6e281cc654a2da9e4f08660c620a1051337fa8965f727eb19190121038262a6c6cec93c2d3ecd6c6072efea86d02ff8e3328bbd0242b20af3425990ac")
	//hex, err := hex.DecodeString("3045022100921cdc9089a5abb4255caaab5a8839a5260842cf57ad3adcb2f7da1e63152cda0220334a96f6b08ece6f746a9cabdc6e4b7d8957503d492a20e6fee310b38fe8d7620102f363c6fb1d71521494da679059a98479d52f7166b941514248a2bfe9bba73173")
	assert.Nil(t, err)
	t.Logf("------:%s", DisasmString(hex))

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

	removed := RemoveOpcode(script, OP_PUSHDATA1)
	assert.NotNil(t, removed)

	removed = RemoveOpcode(removed, OP_PUSHDATA2)
	assert.NotNil(t, removed)

	removed = RemoveOpcode(removed, OP_PUSHDATA4)
	assert.NotNil(t, removed)

	got := DisasmString(removed)
	want := "OP_DUP OP_HASH160 OP_DATA_2 8000 OP_EQUALVERIFY OP_CHECKSIG"
	assert.Equal(t, want, got)
}
