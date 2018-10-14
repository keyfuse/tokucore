// tokucore
//
// Copyright (c) 2018 TokuBlock
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
