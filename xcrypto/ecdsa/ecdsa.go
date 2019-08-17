// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package ecdsa

import (
	"errors"
	"math/big"

	"crypto/ecdsa"
	"crypto/elliptic"
)

// Sign -- generates a deterministic ECDSA signature according to RFC 6979 and BIP62.
func Sign(priv *ecdsa.PrivateKey, hash []byte) (*big.Int, *big.Int, error) {
	c := priv.PublicKey.Curve
	D := priv.D
	N := c.Params().N

	// RFC6979
	k := NonceRFC6979(N, D, hash)

	// point (x1,y1) = k*G
	// r = x1
	r, _ := priv.Curve.ScalarBaseMult(k.Bytes())
	r.Mod(r, N)
	if r.Sign() == 0 {
		return nil, nil, errors.New("calculated R is zero")
	}

	// s = (hash+D*r)/k mod N
	e := HashToInt(c, hash)
	s := new(big.Int).Mul(D, r)
	kinv := new(big.Int).ModInverse(k, N)
	s.Add(s, e).Mul(s, kinv).Mod(s, N)

	// The signature is composed of two values, the r value and the s value.
	// If the s value is greater than N/2, which is not allowed.
	// Just add in some code that if s is greater than N/2, then s = N - s.
	halfOrder := new(big.Int).Rsh(N, 1)
	if s.Cmp(halfOrder) == 1 {
		s.Sub(N, s)
	}
	if s.Sign() == 0 {
		return nil, nil, errors.New("calculated S is zero")
	}

	// Clean.
	k.SetInt64(0)
	return r, s, nil
}

// Verify -- calls ecdsa.Verify to verify the signature of hash using the public key.
// Returns true if the signature is valid, false otherwise.
func Verify(pub *ecdsa.PublicKey, hash []byte, r *big.Int, s *big.Int) bool {
	return ecdsa.Verify(pub, hash, r, s)
}

// HashToInt -- converts a hash value to an integer. There is some disagreement
// about how this is done. [NSA] suggests that this is done in the obvious
// manner, but [SECG] truncates the hash to the bit-length of the curve order
// first. We follow [SECG] because that's what OpenSSL does. Additionally,
// OpenSSL right shifts excess bits from the number if the hash is too large
// and we mirror that too.
// This is borrowed from crypto/ecdsa.
func HashToInt(c elliptic.Curve, hash []byte) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}
