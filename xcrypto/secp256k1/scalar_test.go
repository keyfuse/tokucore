// Copyright 2019 by KeyFuse Labs
// BSD License

package secp256k1

import (
	"math/big"
	"testing"
)

func TestScalar(t *testing.T) {
	s256 := SECP256K1()
	a := big.NewInt(3)
	b := big.NewInt(7)
	s1 := NewScalar(a, b)

	a1 := big.NewInt(2)
	b1 := big.NewInt(4)
	s2 := NewScalar(a1, b1)

	a2 := big.NewInt(4)
	b2, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007908834671653", 10)

	s3 := s1.Add(s256, s2)
	if s3.X.Cmp(a2) != 0 || s3.Y.Cmp(b2) != 0 {
		t.Fatal("not.equal")
	}
}
