// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvVect(t *testing.T) {
	{
		inv := NewInvVect(InvTypeTx, bytes.Repeat([]byte{0x00}, 32))
		typ := inv.Type.String()
		assert.Equal(t, "MSG_TX", typ)
	}
	{
		inv := NewInvVect(255, bytes.Repeat([]byte{0x00}, 32))
		typ := inv.Type.String()
		assert.Equal(t, "Unknown.InvType(255)", typ)
	}
}
