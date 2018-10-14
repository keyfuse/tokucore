// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

func opPushData(vm *Engine) error {
	data := vm.instruction.data
	vm.dstack.PushByteArray(data)
	return nil
}

// opN --
// a common handler for the small integer data push opcodes.
// It pushes the numeric value the opcode represents (which will be from 1 to 16) onto the data stack.
func opN(vm *Engine) error {
	vm.dstack.PushInt(ScriptNum((vm.instruction.op.value - (OP_1 - 1))))
	return nil
}
