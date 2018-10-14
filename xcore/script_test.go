// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScriptClass(t *testing.T) {
	script := []byte{0x01}
	class := GetScriptClass(script)
	assert.Equal(t, NonStandardTy, class)
}
