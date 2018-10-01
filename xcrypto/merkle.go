// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcrypto

import (
	"math"
)

const (
	merkleHashSize = 32
)

// Merkle --
type Merkle struct {
	Hashs [][]byte
}

// NewMerkle -- create new Merkle.
// children nodes.  A diagram depicting how this works for bitcoin transactions
// where h(x) is a double sha256 follows:
//
//	         root = h1234 = h(h12 + h34)
//	        /                           \
//	  h12 = h(h1 + h2)            h34 = h(h3 + h4)
//	   /            \              /            \
//	h1 = h(tx1)  h2 = h(tx2)    h3 = h(tx3)  h4 = h(tx4)
//
// The above stored as a linear array is as follows:
//
// 	[h1 h2 h3 h4 h12 h34 root]
//
// As the above shows, the merkle root is always the last element in the array.
func NewMerkle(hashs [][]byte) *Merkle {
	pow2 := upperPow2(len(hashs))
	arraySize := pow2*2 - 1
	merkles := make([][]byte, arraySize)
	copy(merkles, hashs)

	offset := pow2
	for i := 0; i < arraySize-1; i += 2 {
		switch {
		case merkles[i+1] == nil:
			newHash := merkleNode(merkles[i], merkles[i])
			merkles[offset] = newHash
		default:
			newHash := merkleNode(merkles[i], merkles[i+1])
			merkles[offset] = newHash
		}
		offset++
	}
	return &Merkle{Hashs: merkles}
}

func upperPow2(n int) int {
	if n&(n-1) == 0 {
		return n
	}
	exponent := uint(math.Log2(float64(n))) + 1
	return 1 << exponent // 2^exponent
}

func merkleNode(left []byte, right []byte) []byte {
	var hash [merkleHashSize * 2]byte
	copy(hash[:merkleHashSize], left[:])
	copy(hash[merkleHashSize:], right[:])
	return DoubleSha256(hash[:])
}
