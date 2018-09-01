// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstructionBytes(t *testing.T) {
	data0 := bytes.Repeat([]byte{0x01}, 0xff-1)
	data1 := bytes.Repeat([]byte{0x01}, 0xff+1)
	data2 := bytes.Repeat([]byte{0x01}, 0xffff+1)

	script, err := NewScriptBuilder().AddOp(OP_DUP).AddOp(OP_HASH160).AddInt64(128).AddData(data0).AddData(data1).AddData(data2).AddOp(OP_EQUALVERIFY).AddOp(OP_CHECKSIG).Script()
	assert.Nil(t, err)

	instrs, err := NewScriptReader(script).AllInstructions()
	assert.Nil(t, err)
	for _, instr := range instrs {
		_, err := instr.bytes()
		assert.Nil(t, err)
		assert.Equal(t, instr.OpCode(), instr.op.value)
	}
}
