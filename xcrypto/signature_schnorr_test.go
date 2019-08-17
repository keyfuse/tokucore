// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignatureSchnorrSignerAndVerifer(t *testing.T) {
	msg := DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	key1 := PrvKeyFromBytes([]byte{0x01})
	key2 := PrvKeyFromBytes([]byte{0x02})

	{
		signature, err := SchnorrSign(key1, msg)
		assert.Nil(t, err)

		err = SchnorrVerify(key1.PubKey(), msg, signature)
		assert.Nil(t, err)

		err = SchnorrVerify(key2.PubKey(), msg, signature)
		got := err.Error()
		want := "schnorr.signature.verify.failed"
		assert.Equal(t, want, got)
	}
}
func BenchmarkSignatureSchnorrSigner(b *testing.B) {
	msg := DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	key1 := PrvKeyFromBytes([]byte{0x01})
	for n := 0; n < b.N; n++ {
		_, err := SchnorrSign(key1, msg)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkSignatureSchnorrVerifier(b *testing.B) {
	msg := DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	key1 := PrvKeyFromBytes([]byte{0x01})
	signature, err := SchnorrSign(key1, msg)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		err = SchnorrVerify(key1.PubKey(), msg, signature)
		if err != nil {
			panic(err)
		}
	}
}
