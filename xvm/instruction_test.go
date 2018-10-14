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

func TestInstruction(t *testing.T) {
	datas := [][]byte{
		bytes.Repeat([]byte{0x01}, 0xff-1),
		bytes.Repeat([]byte{0x01}, 0xff+1),
		bytes.Repeat([]byte{0x01}, 0xffff+1),
	}
	script, err := NewScriptBuilder().AddOp(OP_DUP).AddOp(OP_HASH160).AddInt64(128).AddData(datas[0]).AddData(datas[1]).AddData(datas[2]).AddOp(OP_EQUALVERIFY).AddOp(OP_CHECKSIG).Script()
	assert.Nil(t, err)

	instrs, err := NewScriptReader(script).AllInstructions()
	assert.Nil(t, err)
	for _, instr := range instrs {
		_, err := instr.bytes()
		assert.Nil(t, err)
		assert.Equal(t, instr.OpCode(), instr.op.value)
		instr.Data()
	}
}
