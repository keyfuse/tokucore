// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xerror

import (
	"testing"
)

func TestStack(t *testing.T) {
	stack := caller()
	t.Logf("%v", stack.trace())
}
