// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

// Trace -- trace for vm.
type Trace struct {
	Step      uint64
	Executed  string
	Stack     string
	Remaining string
}
