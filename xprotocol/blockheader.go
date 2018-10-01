// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"github.com/tokublock/tokucore/xbase"
	"github.com/tokublock/tokucore/xcrypto"
)

// BlockHeader defines information about a block and is used in the bitcoin
// block (MsgBlock) and headers (MsgHeaders) messages.
type BlockHeader struct {
	// Version of the block.  This is not the same as the protocol version.
	Version uint32

	// Hash of the previous block header in the block chain.
	PrevBlock []byte

	// Merkle tree reference to hash of all transactions for the block.
	MerkleRoot []byte

	// Time the block was created.  This is, unfortunately, encoded as a
	// uint32 on the wire and therefore is limited to 2106.
	Timestamp uint32

	// Difficulty target for the block.
	Bits uint32

	// Nonce used to generate the block.
	Nonce uint32
}

// BlockHash -- calc the block hash.
func (b *BlockHeader) BlockHash() []byte {
	buffer := xbase.NewBuffer()

	buffer.WriteU32(b.Version)
	buffer.WriteBytes(b.PrevBlock)
	buffer.WriteBytes(b.MerkleRoot)
	buffer.WriteU32(b.Timestamp)
	buffer.WriteU32(b.Bits)
	buffer.WriteU32(b.Nonce)
	return xcrypto.DoubleSha256(buffer.Bytes())
}

// Size -- the size of the blockheader.
func (b *BlockHeader) Size() int {
	return 4 + len(b.PrevBlock) + len(b.MerkleRoot) + 4 + 4 + 4
}
