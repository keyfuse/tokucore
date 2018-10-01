// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xerror

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := NewError(Errors, ER_UNKNOWN)
	eStr := err.Error()
	assert.NotNil(t, eStr)
	t.Logf("%+v", err)
	t.Logf("%v", err)
}

func TestErrorNumber(t *testing.T) {
	err := NewError(Errors, -1)
	assert.Equal(t, ER_UNKNOWN, err.Num)
}
