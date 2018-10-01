// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tokublock/tokucore/xbase"
)

func TestBlockHeader(t *testing.T) {
	prev, err := xbase.NewIDFromString("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")
	assert.Nil(t, err)
	merk, err := xbase.NewIDFromString("0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098")
	assert.Nil(t, err)

	ts, err := time.Parse(time.RFC3339, "2009-01-09T02:54:25Z")
	assert.Nil(t, err)

	hdr := &BlockHeader{
		Version:    1,
		PrevBlock:  prev,
		MerkleRoot: merk,
		Timestamp:  uint32(ts.Unix()),
		Bits:       486604799,
		Nonce:      2573394689,
	}

	want := "00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048"
	got := xbase.NewIDToString(hdr.BlockHash())
	assert.Equal(t, want, got)
}
