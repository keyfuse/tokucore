// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"fmt"
	"strings"

	"github.com/tokublock/tokucore/xbase"
	"github.com/tokublock/tokucore/xerror"
)

// ScriptReader -- parsing custom script to parsed opcodes.
type ScriptReader struct {
	scriptLen int
	buffer    *xbase.Buffer
}

// NewScriptReader -- creates new ScriptReader.
func NewScriptReader(script []byte) *ScriptReader {
	return &ScriptReader{
		scriptLen: len(script),
		buffer:    xbase.NewBufferReader(script),
	}
}

// NextInstruction -- returns the current instruction.
func (r *ScriptReader) NextInstruction() (*Instruction, error) {
	buffer := r.buffer
	if buffer.End() {
		return nil, nil
	}

	rop, err := buffer.ReadU8()
	if err != nil {
		return nil, xerror.NewError(Errors, ER_SCRIPT_INSTRUCTION_READ_ERROR, buffer.Len()-buffer.Seek())
	}

	op, ok := opcodes[rop]
	if !ok {
		return nil, xerror.NewError(Errors, ER_SCRIPT_INSTRUCTION_UNKNOWN, rop)
	}
	instr := &Instruction{
		op: &op,
	}

	// Parse the data of the instruction.
	switch {
	// OP_1NEGATE, OP_0, and OP_[1-16]
	// No addtional data, represent the data themselves.
	case op.length == 1:

	// OP_DATA_[1-75]
	// Data pushes of specific lengths.
	case op.length > 1:
		data, err := buffer.ReadBytes(op.length - 1)
		if err != nil {
			return nil, xerror.NewError(Errors, ER_SCRIPT_OPCODE_READ_ERROR, op.name, op.length, buffer.Len()-buffer.Seek())
		}
		instr.data = data

	// OP_PUSHDATAP{1,2,4}
	// Data pushes with parsed lengths.
	case op.length < 0:
		var l uint32
		var err error

		switch op.length {
		case -1:
			l1, err1 := buffer.ReadU8()
			l = uint32(l1)
			err = err1
		case -2:
			l, err = buffer.ReadU16()
		case -4:
			l, err = buffer.ReadU32()
		default:
			return nil, xerror.NewError(Errors, ER_SCRIPT_OPCODE_SIZE_MALFORMED, op.name, -op.length)
		}
		if err != nil {
			return nil, xerror.NewError(Errors, ER_SCRIPT_OPCODE_READ_ERROR, op.name, op.length, buffer.Len()-buffer.Seek())
		}

		// Datas.
		data, err := buffer.ReadBytes(int(l))
		if err != nil {
			return nil, xerror.NewError(Errors, ER_SCRIPT_OPCODE_READ_ERROR, op.name, op.length, buffer.Len()-buffer.Seek())
		}
		instr.data = data
	}
	return instr, nil
}

// AllInstructions -- returns all instructions.
func (r *ScriptReader) AllInstructions() ([]Instruction, error) {
	instrs := make([]Instruction, 0, r.scriptLen)
	for {
		instr, err := r.NextInstruction()
		if err != nil {
			return nil, err
		}
		if instr == nil {
			break
		}
		instrs = append(instrs, *instr)
	}
	return instrs, nil
}

// DisasmRemaining -- disasming the remaining opcodes in the reader.
func (r *ScriptReader) DisasmRemaining() string {
	return DisasmString(r.buffer.Remaining())
}

// DisasmString -- disasming the opcodes to the string instruction.
func DisasmString(script []byte) string {
	instrs, err := NewScriptReader(script).AllInstructions()
	if err != nil {
		return err.Error()
	}

	line := []string{}
	for _, instr := range instrs {
		line = append(line, instr.op.name)
		if instr.data != nil {
			line = append(line, fmt.Sprintf("%x", instr.data))
		}
	}
	return strings.Join(line, " ")
}

// RemoveOpcode -- remove the opcode.
func RemoveOpcode(buf []byte, opcode byte) []byte {
	reader := NewScriptReader(buf)
	instrs, err := reader.AllInstructions()
	if err != nil {
		return buf
	}

	rets := make([]byte, 0, len(buf))
	for _, instr := range instrs {
		if instr.op.value != opcode {
			instrBytes, err := instr.bytes()
			if err != nil {
				return buf
			}
			rets = append(rets, instrBytes...)
		}
	}
	return rets
}
