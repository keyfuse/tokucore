// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"github.com/tokublock/tokucore/xerror"
)

// Error type.
const (
	ER_SCRIPT_INSTRUCTION_UNKNOWN     int = 1101
	ER_SCRIPT_INSTRUCTION_READ_ERROR  int = 1102
	ER_SCRIPT_OPCODE_READ_ERROR       int = 1103
	ER_SCRIPT_OPCODE_SIZE_MALFORMED   int = 1104
	ER_SCRIPT_STACK_INDEX_INVALID     int = 1110
	ER_SCRIPT_STACK_OPERATION_INVALID int = 1111
	ER_SCRIPTNUM_TOO_BIG              int = 1201
	ER_SCRIPTNUM_MINIMAL_DATA         int = 1202
	ER_VM_EXEC_OPCODE_FAILED          int = 1300
)

// Errors -- the jump table of error.
var Errors = map[int]*xerror.Error{
	ER_SCRIPT_INSTRUCTION_UNKNOWN:     {Num: ER_SCRIPT_INSTRUCTION_UNKNOWN, State: "TS000", Message: "script.instruction.unknow[%v]"},
	ER_SCRIPT_INSTRUCTION_READ_ERROR:  {Num: ER_SCRIPT_INSTRUCTION_READ_ERROR, State: "TS000", Message: "script.read.instruction.error.remainning[%v]"},
	ER_SCRIPT_OPCODE_READ_ERROR:       {Num: ER_SCRIPT_OPCODE_READ_ERROR, State: "TS000", Message: "script.read.opcode[%v].requires[%v].bytes.but.remainning[%v]"},
	ER_SCRIPT_OPCODE_SIZE_MALFORMED:   {Num: ER_SCRIPT_OPCODE_SIZE_MALFORMED, State: "TS000", Message: "script.opcode[%v].size[%v].invalid"},
	ER_SCRIPT_STACK_INDEX_INVALID:     {Num: ER_SCRIPT_STACK_INDEX_INVALID, State: "TS000", Message: "script.stack.index[%v].invalid.for.stack.size[%v]"},
	ER_SCRIPT_STACK_OPERATION_INVALID: {Num: ER_SCRIPT_STACK_OPERATION_INVALID, State: "TS000", Message: "script.stack.operation[%v][%v].invalid"},
	ER_SCRIPTNUM_TOO_BIG:              {Num: ER_SCRIPTNUM_TOO_BIG, State: "TS000", Message: "script.num.value.encoded.as[%x].is.[%d]bytes.which.exceeds.the.max.allowed.of.[%d]"},
	ER_SCRIPTNUM_MINIMAL_DATA:         {Num: ER_SCRIPTNUM_MINIMAL_DATA, State: "TS000", Message: "script.num.value.encoded.as[%x].is.not.minimally.encoded"},
	ER_VM_EXEC_OPCODE_FAILED:          {Num: ER_VM_EXEC_OPCODE_FAILED, State: "TVM00", Message: "vm.execute.opcode[%v].failed"},
}
