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

func TestHomomorphismAddPlaintext(t *testing.T) {
	pk, sk, err := GenerateKeyPair(1024)
	assert.Nil(t, err)

	type args struct {
		msg int64
		pt  int64
	}

	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"ct2.input.valid.Add(9)",
			args{2, 9},
			11,
			false,
		},
		{
			"ct36.input.valid.Add(36)",
			args{36, 36},
			72,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := big.NewInt(tt.args.msg)
			pt := big.NewInt(tt.args.pt)

			ct, _ := pk.Encrypt(msg)
			got, err := pk.AddPlaintext(ct, pt)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

			want := big.NewInt(tt.want)
			sum, err := sk.Decrypt(got)
			assert.Nil(t, err)
			assert.Equal(t, want, sum)
			t.Logf("\n homomorphic.raw.number:\t[%v]\n homomorphic.encrypt:\t[%X]\n homomorphic.add(%v):\t[%X]\n homomorphic.decrypt:\t[%v]", tt.args.msg, ct.Bytes(), tt.args.pt, got.Bytes(), sum)
		})
	}
}

func TestPublicKeyHomomorphismAdd(t *testing.T) {
	pk, sk, err := GenerateKeyPair(1024)
	assert.Nil(t, err)
	b2 := new(big.Int).SetInt64(2)
	b245 := new(big.Int).SetInt64(245)
	ct2, _ := pk.Encrypt(b2)
	ct245, _ := pk.Encrypt(b245)

	type args struct {
		ct1 *big.Int
		ct2 *big.Int
	}

	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"input.invalid",
			args{zero, zero},
			0,
			true,
		},
		{
			"inputs.valid",
			args{ct2, ct245},
			247,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct1 := tt.args.ct1
			ct2 := tt.args.ct2

			got, err := pk.Add(ct1, ct2)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

			sum, err := sk.Decrypt(got)
			assert.Nil(t, err)
			want := big.NewInt(tt.want)
			assert.Equal(t, want, sum)
		})
	}
}

func TestHomomorphismMultPlaintext(t *testing.T) {
	pk, sk, _ := GenerateKeyPair(1024)

	b2 := big.NewInt(2)
	b36 := big.NewInt(36)
	ct2, _ := pk.Encrypt(b2)
	ct36, _ := pk.Encrypt(b36)

	type args struct {
		ct *big.Int
		pt int64
	}

	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"input.invalid",
			args{zero, 0},
			0,
			true,
		},
		{
			"input.valid",
			args{ct2, 2},
			4,
			false,
		},
		{
			"input.valid",
			args{ct36, 36},
			1296,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := big.NewInt(tt.args.pt)
			got, err := pk.MultPlaintext(tt.args.ct, pt)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

			sum, err := sk.Decrypt(got)
			assert.Nil(t, err)
			want := big.NewInt(tt.want)
			assert.Equal(t, want, sum)
		})
	}
}

func TestHomomorphismSub(t *testing.T) {
	pk, sk, _ := GenerateKeyPair(1024)

	c2, _ := pk.Encrypt(big.NewInt(2))
	c3, _ := pk.Encrypt(big.NewInt(3))
	c4, _ := pk.Encrypt(big.NewInt(4))
	c5, _ := pk.Encrypt(big.NewInt(5))
	b23578 := big.NewInt(23578)
	c23578, _ := pk.Encrypt(b23578)
	b115 := big.NewInt(115)
	c115, _ := pk.Encrypt(b115)

	tests := []struct {
		name string
		ct1  *big.Int
		ct2  *big.Int
		want int64
	}{
		{"3-2", c3, c2, 3 - 2},
		{"5-2", c5, c2, 5 - 2},
		{"5-3", c5, c3, 5 - 3},
		{"5-4", c5, c4, 5 - 4},
		{"5-4", c23578, c115, 23578 - 115},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pk.Sub(tt.ct1, tt.ct2)

			sum, err := sk.Decrypt(got)
			assert.Nil(t, err)
			want := big.NewInt(tt.want)
			assert.Equal(t, want, sum)
		})
	}
}

func TestHomomorphismDivPlaintext(t *testing.T) {
	pk, sk, _ := GenerateKeyPair(1024)

	b2 := big.NewInt(2)
	ct2, _ := pk.Encrypt(b2)

	b36 := big.NewInt(36)
	ct36, _ := pk.Encrypt(b36)

	type args struct {
		ct *big.Int
		pt int64
	}

	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"input.invalid",
			args{zero, 0},
			0,
			true,
		},
		{
			"input.valid",
			args{ct2, 2},
			1,
			false,
		},
		{
			"input.valid",
			args{ct36, 2},
			18,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := big.NewInt(tt.args.pt)
			got, err := pk.DivPlaintext(tt.args.ct, pt)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

			sum, err := sk.Decrypt(got)
			assert.Nil(t, err)
			want := big.NewInt(tt.want)
			assert.Equal(t, want, sum)
		})
	}
}

func TestHomomorphismAddMultDivTest(t *testing.T) {
	pk, sk, err := GenerateKeyPair(1024)
	assert.Nil(t, err)

	// e(pk1)
	pk1 := big.NewInt(2)
	ct, err := pk.Encrypt(pk1)
	assert.Nil(t, err)

	// r*e(pk1)
	r := big.NewInt(3)
	ct, err = pk.MultPlaintext(ct, r)
	assert.Nil(t, err)

	// r*e(pk1)*pk2
	pk2 := big.NewInt(4)
	ct, err = pk.MultPlaintext(ct, pk2)
	assert.Nil(t, err)

	// z+r*e(pk1)*pk2
	z := big.NewInt(11)
	ct, err = pk.AddPlaintext(ct, z)
	assert.Nil(t, err)

	// (z+r*e(pk1)*pk2)/k2
	k2 := big.NewInt(5)
	ct, err = pk.DivPlaintext(ct, k2)
	assert.Nil(t, err)

	want := big.NewInt(7)
	got, err := sk.Decrypt(ct)
	assert.Equal(t, want, got)
}
