// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcrypto

import (
	"errors"
	"math/big"

	"crypto/ecdsa"
	"crypto/elliptic"
)

// EcdsaSign -- generates a deterministic ECDSA signature according to RFC 6979 and BIP62.
func EcdsaSign(priv *ecdsa.PrivateKey, hash []byte) (*big.Int, *big.Int, error) {
	c := priv.PublicKey.Curve
	D := priv.D
	N := c.Params().N

	// RFC6979
	k := nonceRFC6979(N, D, hash)

	// point (x1,y1) = k*G
	// r = x1
	r, _ := priv.Curve.ScalarBaseMult(k.Bytes())
	r.Mod(r, N)
	if r.Sign() == 0 {
		return nil, nil, errors.New("calculated R is zero")
	}

	// s = (k^-1 (hash + D * r) mod N
	halfOrder := new(big.Int).Rsh(N, 1)
	kinv := new(big.Int).ModInverse(k, N)
	e := hashToInt(hash, c)
	s := new(big.Int).Mul(D, r)
	s.Add(s, e).Mul(s, kinv).Mod(s, N)

	// https://bitcoin.stackexchange.com/questions/68254/how-can-i-fix-this-non-canonical-signature-s-value-is-unnecessarily-high?rq=1
	// The signature is composed of two values, the r value and the s value.
	// If the s value is greater than N/2, which is not allowed.
	// Just add in some code that if s is greater than N/2, then s = N - s.
	if s.Cmp(halfOrder) == 1 {
		s.Sub(N, s)
	}
	if s.Sign() == 0 {
		return nil, nil, errors.New("calculated S is zero")
	}
	return r, s, nil
}

// EcdsaVerify -- calls ecdsa.Verify to verify the signature of hash using the public key.
// Returns true if the signature is valid, false otherwise.
func EcdsaVerify(pub *ecdsa.PublicKey, hash []byte, r *big.Int, s *big.Int) bool {
	return ecdsa.Verify(pub, hash, r, s)
}

// hashToInt converts a hash value to an integer. There is some disagreement
// about how this is done. [NSA] suggests that this is done in the obvious
// manner, but [SECG] truncates the hash to the bit-length of the curve order
// first. We follow [SECG] because that's what OpenSSL does. Additionally,
// OpenSSL right shifts excess bits from the number if the hash is too large
// and we mirror that too.
// This is borrowed from crypto/ecdsa.
func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
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
