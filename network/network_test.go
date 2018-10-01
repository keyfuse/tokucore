// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetwork(t *testing.T) {
	net := TestNet
	net.SetMagic([]byte{0x01}).
		SetPort(1333).
		SetLastBlock(100)

	assert.Equal(t, []byte{0x01}, net.Magic)
	assert.Equal(t, uint32(1333), net.Port)
	assert.Equal(t, uint32(100), net.LastBlock)
}
