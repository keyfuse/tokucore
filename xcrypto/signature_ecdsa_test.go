// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignatureEcdsaSignerAndVerifer(t *testing.T) {
	msg := DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	key1 := PrvKeyFromBytes([]byte{0x01})
	key2 := PrvKeyFromBytes([]byte{0x02})

	{
		signature, err := EcdsaSign(key1, msg)
		assert.Nil(t, err)

		err = EcdsaVerify(key1.PubKey(), msg, signature)
		assert.Nil(t, err)

		err = EcdsaVerify(key2.PubKey(), msg, signature)
		got := err.Error()
		want := "ecdsa.signature.verify.failed"
		assert.Equal(t, want, got)
	}
}

func BenchmarkSignatureEcdsaSigner(b *testing.B) {
	msg := DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	key1 := PrvKeyFromBytes([]byte{0x01})
	for n := 0; n < b.N; n++ {
		_, err := EcdsaSign(key1, msg)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkSignatureEcdsaVerifier(b *testing.B) {
	msg := DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	key1 := PrvKeyFromBytes([]byte{0x01})
	signature, err := EcdsaSign(key1, msg)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		err = EcdsaVerify(key1.PubKey(), msg, signature)
		if err != nil {
			panic(err)
		}
	}
}
