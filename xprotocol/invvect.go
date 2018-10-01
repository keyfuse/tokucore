// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"fmt"
)

const (
	// MaxInvPerMsg -- the maximum number of inventory vectors that can be in a
	// single bitcoin inv message.
	MaxInvPerMsg = 50000
)

// InvType -- the allowed types of inventory vectors.  See InvVect.
type InvType uint32

//
const (
	InvTypeError         InvType = 0
	InvTypeTx            InvType = 1
	InvTypeBlock         InvType = 2
	InvTypeFilteredBlock InvType = 3
)

//
var ivStrings = map[InvType]string{
	InvTypeError:         "ERROR",
	InvTypeTx:            "MSG_TX",
	InvTypeBlock:         "MSG_BLOCK",
	InvTypeFilteredBlock: "MSG_FILTERED_BLOCK",
}

// String -- returns the InvType in human-readable form.
func (invtype InvType) String() string {
	if s, ok := ivStrings[invtype]; ok {
		return s
	}
	return fmt.Sprintf("Unknown.InvType(%d)", uint32(invtype))
}

// InvVect -- defines a bitcoin inventory vector which is used to describe data.
type InvVect struct {
	Type InvType // Type of data
	Hash []byte  // Hash of the data
}

// NewInvVect -- returns a new InvVect using the provided type and hash.
func NewInvVect(typ InvType, hash []byte) *InvVect {
	return &InvVect{
		Type: typ,
		Hash: hash,
	}
}

// Size -- size of the InvVect.
func (inv *InvVect) Size() int {
	return 4 + len(inv.Hash)
}
