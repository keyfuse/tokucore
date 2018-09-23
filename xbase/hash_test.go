// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xbase

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	bytes := bytes.Repeat([]byte{0x02}, 32)
	want := NewIDToString(bytes)
	got, err := NewIDFromString(want)
	assert.Nil(t, err)
	assert.Equal(t, got, bytes)
}
