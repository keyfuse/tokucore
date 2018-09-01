// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"bytes"

	"github.com/tokublock/tokucore/xerror"
)

// opcodeFalse --
// pushes an empty array to the data stack to represent false.  Note
// that 0, when encoded as a number according to the numeric encoding consensus rules, is an empty array.
func opFalse(vm *Engine) error {
	vm.dstack.PushByteArray(nil)
	return nil
}

// opEqual --
// removes the top 2 items of the data stack, compares them as raw
// bytes, and pushes the result, encoded as a boolean, back to the stack.
//
// Stack : [... x1 x2] -> [... bool]
func opEqual(vm *Engine) error {
	a, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}
	b, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}
	vm.dstack.PushBool(bytes.Equal(a, b))
	return nil
}

func opEqualVerify(vm *Engine) error {
	err := opEqual(vm)
	if err != nil {
		return err
	}
	return equalVerify(vm, xerror.NewError(Errors, ER_VM_EXEC_OPCODE_FAILED, "opEqualVerify"))
}

func equalVerify(vm *Engine, operr error) error {
	verified, err := vm.dstack.PopBool()
	if err != nil {
		return err
	}
	if !verified {
		return operr
	}
	return nil
}
