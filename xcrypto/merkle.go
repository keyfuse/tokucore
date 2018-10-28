// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcrypto

import (
	"bytes"
	"math"
)

const (
	merkleHashSize = 32
	positionLeft   = "left"
	positionRight  = "right"
)

// Node -- node for prove.
type Node struct {
	Hash     []byte
	Parent   []byte
	Position string
}

// Merkle --
type Merkle struct {
	levels [][][]byte
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
// If the right is nil(the size is not power of 2), hashed with itself.
func NewMerkle(hashs [][]byte) *Merkle {
	n := len(hashs)
	depth := int(math.Ceil(math.Log2(float64(n))))
	pow2 := int(math.Pow(2, float64(depth)))
	levels := make([][][]byte, depth+1)
	levels[depth] = make([][]byte, pow2)
	copy(levels[depth], hashs)

	// Build the merkle tree with levels.
	for i := depth; i > 0; i-- {
		offset := 0
		level := levels[i]
		next := make([][]byte, len(level)/2)
		for j := 0; j < len(level); j += 2 {
			left := level[j]
			right := level[j+1]
			if right == nil {
				next[offset] = merkleNode(left, left)
				break
			} else {
				next[offset] = merkleNode(left, right)
			}
			offset++
		}
		levels[i-1] = next
	}
	return &Merkle{levels: levels}
}

func merkleNode(left []byte, right []byte) []byte {
	var hash [merkleHashSize * 2]byte
	copy(hash[:merkleHashSize], left[:])
	copy(hash[merkleHashSize:], right[:])
	return DoubleSha256(hash[:])
}

// Root -- returns the merkle root.
func (m *Merkle) Root() []byte {
	return m.levels[0][0]
}

// Proofs -- gets the proof path for this leaf.
func (m *Merkle) Proofs(leaf []byte) []Node {
	var path []Node

	index := -1
	depth := len(m.levels) - 1
	for i, le := range m.levels[depth] {
		if bytes.Equal(leaf, le) {
			index = i
			break
		}
	}

	if index > -1 {
		for i := depth; i > 0; i-- {
			var parent []byte
			var siblingHash []byte
			var siblingPostion string

			level := m.levels[i]
			if (index % 2) != 0 {
				siblingHash = level[index-1]
				siblingPostion = positionLeft
				parent = merkleNode(siblingHash, level[index])
			} else {
				siblingHash = level[index+1]
				if siblingHash == nil {
					siblingHash = level[index]
				}
				siblingPostion = positionRight
				parent = merkleNode(level[index], siblingHash)
			}
			index = index / 2
			path = append(path, Node{Hash: siblingHash, Parent: parent, Position: siblingPostion})
		}
	}
	return path
}

// Verify -- used to verify the leaf contained in the merkle tree.
func (m *Merkle) Verify(leaf []byte, root []byte, path []Node) bool {
	hash := leaf
	for _, node := range path {
		switch node.Position {
		case positionLeft:
			hash = merkleNode(node.Hash, hash)
		case positionRight:
			if node.Hash == nil {
				hash = merkleNode(hash, hash)
			} else {
				hash = merkleNode(hash, node.Hash)
			}
		default:
			return false
		}
		if !bytes.Equal(node.Parent, hash) {
			return false
		}
	}
	return bytes.Equal(hash, root)
}
