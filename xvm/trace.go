// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xvm

// Trace -- trace for vm.
type Trace struct {
	Step      uint64
	Executed  string
	Stack     string
	Remaining string
}
