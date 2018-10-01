// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xerror

import (
	"testing"
)

func TestStack(t *testing.T) {
	stack := caller()
	t.Logf("%v", stack.trace())
}
