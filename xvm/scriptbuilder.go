// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"encoding/hex"

	"github.com/tokublock/tokucore/xbase"
	"github.com/tokublock/tokucore/xcrypto"
)

// ScriptBuilder -- for building custom scripts
//
// For example, the following would build a 2-of-3 multisig script for usage in
// a pay-to-script-hash (although in this situation MultiSigScript() would be a
// better choice to generate the script):
// 	builder := NewScriptBuilder()
// 	builder.AddOp(OP_2).AddData(pubKey1).AddData(pubKey2)
// 	builder.AddData(pubKey3).AddOp(OP_3)
// 	builder.AddOp(OP_CHECKMULTISIG)
// 	script, err := builder.Script()
type ScriptBuilder struct {
	buffer *xbase.Buffer
	err    error
}

// NewScriptBuilder -- creates ScriptBuilder
func NewScriptBuilder() *ScriptBuilder {
	return &ScriptBuilder{
		buffer: xbase.NewBuffer(),
	}
}

// AddData -- push the passed data to the end of the script
func (b *ScriptBuilder) AddData(data []byte) *ScriptBuilder {
	l := len(data)
	buffer := b.buffer

	switch {
	case l == 0:
		buffer.WriteU8(OP_0)
		return b

	// OP_1NEGATE, OP_0, and OP_[1-16]
	case l == 1:
		switch {
		// OP_0
		case data[0] == 0:
			buffer.WriteU8(OP_0)
			return b
		// OP_[1-16]
		case data[0] <= 16:
			buffer.WriteU8((OP_1 - 1) + data[0])
			return b
		// OP_1NEGATE
		case data[0] == 0x81:
			buffer.WriteU8(OP_1NEGATE)
			return b
		}
		fallthrough

	// OP_DATA_[1-75]
	case l < OP_PUSHDATA1:
		buffer.WriteU8(uint8((OP_DATA_1 - 1) + l))

	// OP_PUSHDATAP1
	case l <= 0xff:
		buffer.WriteU8(OP_PUSHDATA1)
		buffer.WriteU8(uint8(l))

	// OP_PUSHDATAP2
	case l <= 0xffff:
		buffer.WriteU8(OP_PUSHDATA2)
		buffer.WriteU16(uint32(l))

	// OP_PUSHDATAP4
	default:
		buffer.WriteU8(OP_PUSHDATA4)
		buffer.WriteU32(uint32(l))
	}
	buffer.WriteBytes(data)
	return b
}

// AddOp -- push the passed opcode to the end of script
func (b *ScriptBuilder) AddOp(op byte) *ScriptBuilder {
	b.buffer.WriteU8(op)
	return b
}

// AddInt64 -- push the passed int64 to the end of script
func (b *ScriptBuilder) AddInt64(val int64) *ScriptBuilder {
	switch {
	case val == 0:
		b.buffer.WriteU8(OP_0)
	case val == -1 || (val >= 1 && val <= 16):
		b.buffer.WriteU8(byte(OP_1 + val - 1))
	default:
		b.AddData(ScriptNum(val).Bytes())
	}
	return b
}

// Reset -- resets the script so it has no content
func (b *ScriptBuilder) Reset() *ScriptBuilder {
	b.buffer.Reset()
	b.err = nil
	return b
}

// Script -- returns the currently built script
// When any errors occurred while building the script,
// the script will be returned up the point of the first error along with the error.
func (b *ScriptBuilder) Script() ([]byte, error) {
	return b.buffer.Bytes(), b.err
}

// Hash160 -- returns the hash160 of the script bytes.
func (b *ScriptBuilder) Hash160() ([]byte, error) {
	return xcrypto.Hash160(b.buffer.Bytes()), b.err
}

// parseHex -- parse hex string into a []byte.
func parseHex(tok string) ([]byte, error) {
	if !strings.HasPrefix(tok, "0x") {
		return nil, errors.New("not.a.hex.number")
	}
	return hex.DecodeString(tok[2:])
}

// Load -- load the string script and parse to instructions.
func (b *ScriptBuilder) Load(script string) *ScriptBuilder {
	script = strings.Replace(script, "\n", " ", -1)
	script = strings.Replace(script, "\t", " ", -1)
	tokens := strings.Split(script, " ")
	for i, token := range tokens {
		op, ok := opcodesByName[token]
		if ok {
			b.AddOp(op)
		} else {
			if num, err := strconv.ParseInt(token, 10, 64); err == nil { // Num.
				b.AddInt64(num)
			} else if bts, err := parseHex(token); err == nil { // Hex.
				b.buffer.WriteBytes(bts)
			} else { // Data
				tlen := len(token)
				if tlen > 2 && token[0] == '\'' && token[tlen-1] == '\'' {
					datas, _ := hex.DecodeString(token[1 : tlen-1])
					b.AddData(datas)
				} else {
					b.err = fmt.Errorf("unknow.token[%v].at:%v", token, i)
					return b
				}
			}
		}
	}
	return b
}
