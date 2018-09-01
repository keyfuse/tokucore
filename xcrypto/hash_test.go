// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcrypto

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name string
		fn   func([]byte) []byte
		hash string
	}{
		{
			name: "Sha256",
			fn:   Sha256,
			hash: "6a323ef9e4e70ff872dcd0ef1e78efee84471517393e93ddb8306a08c93032f7",
		},
		{
			name: "DoubleSha256",
			fn:   DoubleSha256,
			hash: "6dc2d42e583162dac8ef92d5c305f914951e103c89b7d58a27c5572178fa5bf5",
		},

		{
			name: "Ripemd160",
			fn:   Ripemd160,
			hash: "e5b38131874755792ffa2af4d8400a8ea4f8dac4",
		},
		{
			name: "Hash160",
			fn:   Hash160,
			hash: "e624333290d16e4bef678ec2eba70063f70544e2",
		},
	}

	for _, test := range tests {
		hash := fmt.Sprintf("%x", test.fn([]byte(test.name)))
		assert.Equal(t, test.hash, hash)
	}
}

func TestHashCommon(t *testing.T) {
	// Ripemd160Size.
	size := Ripemd160Size()
	assert.Equal(t, 20, size)

	// BytesToBigInt.
	bint := BytesToBigInt([]byte{0x01})
	assert.Equal(t, new(big.Int).SetBytes([]byte{0x01}), bint)
}
