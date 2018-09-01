// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

// opAdd --
// treats the top two items on the data stack as integers and replaces them with their sum.
func opAdd(vm *Engine) error {
	// v0.
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	// v1.
	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	// Add.
	vm.dstack.PushInt(v0 + v1)
	return nil
}
