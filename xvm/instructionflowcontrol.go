// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"fmt"

	"github.com/tokublock/tokucore/xerror"
)

// opReturn --
func opReturn(vm *Engine) error {
	return nil
}

// opIf -- treats the top item on the data stack as a boolean and removes it.
//
// <expression> if [statements] [else [statements]] endif
//
// Stack:
// [... bool] -> [...]
// Conditional stack transformation: [...] -> [... OpCondValue]
func opIf(vm *Engine) error {
	condVal := OpCondFalse
	if !vm.branchShouldSkip() {
		ok, err := popIfBool(vm)
		if err != nil {
			return err
		}
		if ok {
			condVal = OpCondTrue
		}
	} else {
		condVal = OpCondSkip
	}
	vm.cstack = append(vm.cstack, condVal)
	return nil
}

// opNotIf -- treats the top item on the data stack as a boolean and removes it.
//
// <expression> notif [statements] [else [statements]] endif
//
// Stack:
// [... bool] -> [...]
// Conditional stack transformation: [...] -> [... OpCondValue]
func opNotIf(vm *Engine) error {
	condVal := OpCondFalse
	if !vm.branchShouldSkip() {
		ok, err := popIfBool(vm)
		if err != nil {
			return err
		}
		if !ok {
			condVal = OpCondTrue
		}
	} else {
		condVal = OpCondSkip
	}
	vm.cstack = append(vm.cstack, condVal)
	return nil
}

// opElse -- inverts conditional execution for other half of if/else/endif.
// An error is returned if there has not already been a matching OP_IF.
// Conditional stack transformation: [... OpCondValue] -> [... !OpCondValue]
func opElse(vm *Engine) error {
	if len(vm.cstack) == 0 {
		return xerror.NewError(Errors, ER_VM_EXEC_OPCODE_FAILED, "opElse.no.matching.opcode.to.begin.conditional.execution")
	}
	condIdx := len(vm.cstack) - 1
	switch vm.cstack[condIdx] {
	case OpCondTrue:
		vm.cstack[condIdx] = OpCondFalse
	case OpCondFalse:
		vm.cstack[condIdx] = OpCondTrue
	case OpCondSkip:
	}
	return nil
}

// opEndIf -- terminates a conditional block, removing the value from the conditional execution stack.
// An error is returned if there has not already been a matching OP_IF.
// Conditional stack transformation: [... OpCondValue] -> [...]
func opEndIf(vm *Engine) error {
	if len(vm.cstack) == 0 {
		return xerror.NewError(Errors, ER_VM_EXEC_OPCODE_FAILED, "opEndIf.no.matching.opcode.to.begin.conditional.execution")
	}
	vm.cstack = vm.cstack[:len(vm.cstack)-1]
	return nil
}

// popIfBool --
// require the following: for OP_IF and OP_NOT_IF,
// the top stack item MUST either be an empty byte slice, or [0x01].
// Otherwise, the item at the top of the stack will be popped and interpreted as a boolean.
func popIfBool(vm *Engine) (bool, error) {
	so, err := vm.dstack.PopByteArray()
	if err != nil {
		return false, err
	}

	// The top element MUST have a length of at least one.
	if len(so) > 1 {
		str := fmt.Sprintf("popIfBool.minimal.if.is.active.top.MUST.have.a.length.instead.of:%v", len(so))
		return false, xerror.NewError(Errors, ER_VM_EXEC_OPCODE_FAILED, str)
	}

	// Additionally, if the length is one, then the value MUST be 0x01.
	if len(so) == 1 && so[0] != 0x01 {
		str := fmt.Sprintf("popIfBool.minimal.if.is.active.top.MUST.be.empty.or.0x01.instead:%v", so[0])
		return false, xerror.NewError(Errors, ER_VM_EXEC_OPCODE_FAILED, str)
	}
	return asBool(so), nil
}
