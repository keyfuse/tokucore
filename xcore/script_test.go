// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScript(t *testing.T) {
	script := []byte{0x01}
	_, err := ParseLockingScript(script)
	assert.NotNil(t, err)
}
