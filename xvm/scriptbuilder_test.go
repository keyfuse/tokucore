// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScriptBuilderAddData(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected []byte
	}{
		// BIP0062: Pushing an empty byte sequence must use OP_0.
		{name: "push empty byte sequence", data: nil, expected: []byte{OP_0}},
		{name: "push 1 byte 0x00", data: []byte{0x00}, expected: []byte{OP_0}},

		// BIP0062: Pushing a 1-byte sequence of byte 0x01 through 0x10 must use OP_n.
		{name: "push 1 byte 0x01", data: []byte{0x01}, expected: []byte{OP_1}},
		{name: "push 1 byte 0x02", data: []byte{0x02}, expected: []byte{OP_2}},
		{name: "push 1 byte 0x03", data: []byte{0x03}, expected: []byte{OP_3}},
		{name: "push 1 byte 0x04", data: []byte{0x04}, expected: []byte{OP_4}},
		{name: "push 1 byte 0x05", data: []byte{0x05}, expected: []byte{OP_5}},
		{name: "push 1 byte 0x06", data: []byte{0x06}, expected: []byte{OP_6}},
		{name: "push 1 byte 0x07", data: []byte{0x07}, expected: []byte{OP_7}},
		{name: "push 1 byte 0x08", data: []byte{0x08}, expected: []byte{OP_8}},
		{name: "push 1 byte 0x09", data: []byte{0x09}, expected: []byte{OP_9}},
		{name: "push 1 byte 0x0a", data: []byte{0x0a}, expected: []byte{OP_10}},
		{name: "push 1 byte 0x0b", data: []byte{0x0b}, expected: []byte{OP_11}},
		{name: "push 1 byte 0x0c", data: []byte{0x0c}, expected: []byte{OP_12}},
		{name: "push 1 byte 0x0d", data: []byte{0x0d}, expected: []byte{OP_13}},
		{name: "push 1 byte 0x0e", data: []byte{0x0e}, expected: []byte{OP_14}},
		{name: "push 1 byte 0x0f", data: []byte{0x0f}, expected: []byte{OP_15}},
		{name: "push 1 byte 0x10", data: []byte{0x10}, expected: []byte{OP_16}},
	}

	builder := NewScriptBuilder()
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		builder.Reset().AddData(test.data)
		result, _ := builder.Script()
		if !bytes.Equal(result, test.expected) {
			t.Errorf("ScriptBuilder.AddData #%d (%s) wrong result\n"+
				"got: %x\nwant: %x", i, test.name, result,
				test.expected)
			continue
		}
	}
}

// TestScriptBuilderAddInt64 tests that pushing signed integers to a script via
// the ScriptBuilder API works as expected.
func TestScriptBuilderAddInt64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		val      int64
		expected []byte
	}{
		{name: "push -1", val: -1, expected: []byte{OP_1NEGATE}},
		{name: "push small int 0", val: 0, expected: []byte{OP_0}},
		{name: "push small int 1", val: 1, expected: []byte{OP_1}},
		{name: "push small int 2", val: 2, expected: []byte{OP_2}},
		{name: "push small int 3", val: 3, expected: []byte{OP_3}},
		{name: "push small int 4", val: 4, expected: []byte{OP_4}},
		{name: "push small int 5", val: 5, expected: []byte{OP_5}},
		{name: "push small int 6", val: 6, expected: []byte{OP_6}},
		{name: "push small int 7", val: 7, expected: []byte{OP_7}},
		{name: "push small int 8", val: 8, expected: []byte{OP_8}},
		{name: "push small int 9", val: 9, expected: []byte{OP_9}},
		{name: "push small int 10", val: 10, expected: []byte{OP_10}},
		{name: "push small int 11", val: 11, expected: []byte{OP_11}},
		{name: "push small int 12", val: 12, expected: []byte{OP_12}},
		{name: "push small int 13", val: 13, expected: []byte{OP_13}},
		{name: "push small int 14", val: 14, expected: []byte{OP_14}},
		{name: "push small int 15", val: 15, expected: []byte{OP_15}},
		{name: "push small int 16", val: 16, expected: []byte{OP_16}},
		{name: "push 17", val: 17, expected: []byte{OP_DATA_1, 0x11}},
		{name: "push 65", val: 65, expected: []byte{OP_DATA_1, 0x41}},
		{name: "push 127", val: 127, expected: []byte{OP_DATA_1, 0x7f}},
		{name: "push 128", val: 128, expected: []byte{OP_DATA_2, 0x80, 0}},
		{name: "push 255", val: 255, expected: []byte{OP_DATA_2, 0xff, 0}},
		{name: "push 256", val: 256, expected: []byte{OP_DATA_2, 0, 0x01}},
		{name: "push 32767", val: 32767, expected: []byte{OP_DATA_2, 0xff, 0x7f}},
		{name: "push 32768", val: 32768, expected: []byte{OP_DATA_3, 0, 0x80, 0}},
		{name: "push -2", val: -2, expected: []byte{OP_DATA_1, 0x82}},
		{name: "push -3", val: -3, expected: []byte{OP_DATA_1, 0x83}},
		{name: "push -4", val: -4, expected: []byte{OP_DATA_1, 0x84}},
		{name: "push -5", val: -5, expected: []byte{OP_DATA_1, 0x85}},
		{name: "push -17", val: -17, expected: []byte{OP_DATA_1, 0x91}},
		{name: "push -65", val: -65, expected: []byte{OP_DATA_1, 0xc1}},
		{name: "push -127", val: -127, expected: []byte{OP_DATA_1, 0xff}},
		{name: "push -128", val: -128, expected: []byte{OP_DATA_2, 0x80, 0x80}},
		{name: "push -255", val: -255, expected: []byte{OP_DATA_2, 0xff, 0x80}},
		{name: "push -256", val: -256, expected: []byte{OP_DATA_2, 0x00, 0x81}},
		{name: "push -32767", val: -32767, expected: []byte{OP_DATA_2, 0xff, 0xff}},
		{name: "push -32768", val: -32768, expected: []byte{OP_DATA_3, 0x00, 0x80, 0x80}},
	}

	builder := NewScriptBuilder()
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		builder.Reset().AddInt64(test.val)
		result, err := builder.Script()
		if err != nil {
			t.Errorf("ScriptBuilder.AddInt64 #%d (%s) unexpected "+
				"error: %v", i, test.name, err)
		}
		if !bytes.Equal(result, test.expected) {
			t.Errorf("ScriptBuilder.AddInt64 #%d (%s) wrong result\n"+
				"got: %x\nwant: %x", i, test.name, result,
				test.expected)
			panic(1)
		}
		_, err = builder.Hash160()
		assert.Nil(t, err)
	}
}

func TestScriptBuilderLoadProgram(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		program  string
		expected string
		err      error
	}{
		{
			"ok",
			"OP_1 OP_15 OP_ADD OP_16 OP_EQUAL",
			"OP_1 OP_15 OP_ADD OP_16 OP_EQUAL",
			nil,
		},
		{
			"ok",
			"OP_2 OP_DATA_33 0x03d9914d1e64ad12d673fd9d749d12ca21dc0e9a4648bf04b26e70989027f7cee5 OP_DATA_33 0x027ace0f7d1bef4651cdcf342a8ef2e5e4e031f7ca2c3957cae07a14d2a907933c OP_2 OP_CHECKMULTISIG",
			"OP_2 OP_DATA_33 03d9914d1e64ad12d673fd9d749d12ca21dc0e9a4648bf04b26e70989027f7cee5 OP_DATA_33 027ace0f7d1bef4651cdcf342a8ef2e5e4e031f7ca2c3957cae07a14d2a907933c OP_2 OP_CHECKMULTISIG",
			nil,
		},
		{
			"err",
			"OP_1 OP_15 X OP_ADD OP_16 OP_EQUAL",
			"",
			fmt.Errorf("token.error"),
		},
	}

	builder := NewScriptBuilder()
	for _, test := range tests {
		builder.Reset().Load(test.program)
		script, err := builder.Script()
		if test.err != nil {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			asm := DisasmString(script)
			assert.Equal(t, test.expected, asm)
		}
	}
}
