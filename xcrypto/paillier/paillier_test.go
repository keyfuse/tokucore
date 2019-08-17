// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package paillier

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaillier(t *testing.T) {
	tests := []struct {
		name    string
		bitlen  int
		msg     int64
		wantErr bool
	}{
		{"ok", 1024, 2019, false},
		{"plaintext.invalid", 1024, -1, true},
	}

	for _, test := range tests {
		pk, sk, err := GenerateKeyPair(test.bitlen)
		assert.Nil(t, err)

		pt := big.NewInt(test.msg)
		ct, err := pk.Encrypt(pt)
		if test.wantErr {
			assert.NotNil(t, err)
			continue
		}

		got, err := sk.Decrypt(ct)
		assert.Nil(t, err)
		want := big.NewInt(test.msg)
		assert.Equal(t, want, got)
	}
}

func benchmarkKey(size int, b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		GenerateKeyPair(size)
	}
}

//  go test -v ./xcrypto/paillier -bench=.  -benchtime=1s
func BenchmarkKey1024(b *testing.B) { benchmarkKey(1024, b) }
func BenchmarkKey2048(b *testing.B) { benchmarkKey(2048, b) }
func BenchmarkKey3072(b *testing.B) { benchmarkKey(3072, b) }
func BenchmarkKey4096(b *testing.B) { benchmarkKey(4096, b) }
