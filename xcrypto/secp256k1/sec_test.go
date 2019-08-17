// Copyright 2019 by KeyFuse Labs
// BSD License

package secp256k1

import (
	"testing"
)

func TestSec(t *testing.T) {
	s256 := SECP256K1()
	a, b := s256.ScalarBaseMult([]byte{0x01})
	se := SecMarshal(s256, a, b)
	x, y := SecUnmarshal(s256, se)
	if x.Cmp(a) != 0 || y.Cmp(b) != 0 {
		t.Fatal("not.equal")
	}
}
