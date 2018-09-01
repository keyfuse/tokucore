// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcrypto

import (
	"crypto/sha256"
	"hash"
	"math/big"

	"github.com/tokublock/tokucore/xcrypto/ripemd160"
)

func calcHash(buf []byte, hasher hash.Hash) []byte {
	hasher.Write(buf)
	return hasher.Sum(nil)
}

// Hash160 -- returns ripemd160(sha256) bytes.
func Hash160(data []byte) []byte {
	return Ripemd160(Sha256(data))
}

// Ripemd160 -- returns ripemd160 bytes.
func Ripemd160(data []byte) []byte {
	return calcHash(data, ripemd160.New())
}

// Ripemd160Size -- Size of the ripemd160.
func Ripemd160Size() int {
	return ripemd160.Size
}

// Sha256 -- returns sha256 bytes.
func Sha256(data []byte) []byte {
	return calcHash(data, sha256.New())
}

// DoubleSha256 -- returns sha256(sha256) bytes.
func DoubleSha256(data []byte) []byte {
	hash := calcHash(data, sha256.New())
	return calcHash(hash, sha256.New())
}

// BytesToBigInt -- returns big int for the b bytes represents.
func BytesToBigInt(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}
