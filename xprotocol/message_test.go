// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	commands := []string{
		CommandVersion,
		CommandVersionAck,
		CommandInventory,
		CommandGetData,
		CommandGetHeaders,
		CommandHeaders,
		CommandBlock,
		CommandTx,
		CommandReject,
	}

	for _, cmd := range commands {
		msg := makeEmptyMessage(cmd)
		assert.NotNil(t, msg)
	}
}
