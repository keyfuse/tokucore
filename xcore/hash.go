// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"encoding/hex"
)

// NewTxIDFromString -- returns byte-reverse hash bytes.
func NewTxIDFromString(hashStr string) ([]byte, error) {
	hash, err := hex.DecodeString(hashStr)
	if err != nil {
		return nil, err
	}

	size := len(hash)
	for i := 0; i < size/2; i++ {
		hash[i], hash[size-1-i] = hash[size-1-i], hash[i]
	}
	return hash, nil
}

// NewTxIDToString -- returns byte-reverse hash hex.
func NewTxIDToString(hash []byte) string {
	size := len(hash)
	clone := make([]byte, size)
	copy(clone, hash[:])
	for i := 0; i < size/2; i++ {
		clone[i], clone[size-1-i] = clone[size-1-i], clone[i]
	}
	return hex.EncodeToString(clone[:])
}
