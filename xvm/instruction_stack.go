// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xvm

func opDup(vm *Engine) error {
	return vm.dstack.DupN(1)
}

func op2Dup(vm *Engine) error {
	return vm.dstack.DupN(2)
}

func op3Dup(vm *Engine) error {
	return vm.dstack.DupN(3)
}

// opSwap -- swaps the top two items on the stack.
// Stack:
// [... x1 x2] -> [... x2 x1]
func opSwap(vm *Engine) error {
	return vm.dstack.SwapN(1)
}

// opSize -- pushes the size of the top item of the data stack onto the data stack.
// Stack:
// [... x1] -> [... x1 len(x1)]
func opSize(vm *Engine) error {
	so, err := vm.dstack.PeekByteArray(0)
	if err != nil {
		return err
	}

	vm.dstack.PushInt(ScriptNum(len(so)))
	return nil
}

// opDrop -- removes the top item from the data stack.
// Stack:
// [... x1 x2 x3] -> [... x1 x2]
func opDrop(vm *Engine) error {
	return vm.dstack.DropN(1)
}
