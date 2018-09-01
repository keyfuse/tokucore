// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"hash"
	"math/big"
)

// EcdsaSign -- generates a deterministic ECDSA signature according to RFC 6979 and BIP62.
func EcdsaSign(priv *ecdsa.PrivateKey, hash []byte, alg func() hash.Hash) (r, s *big.Int, err error) {
	c := priv.PublicKey.Curve
	N := c.Params().N
	halfOrder := new(big.Int).Rsh(N, 1)

	generateSecret(N, priv.D, alg, hash, func(k *big.Int) bool {
		inv := new(big.Int).ModInverse(k, N)
		r, _ = priv.Curve.ScalarBaseMult(k.Bytes())
		if r.Cmp(N) == 1 {
			r.Sub(r, N)
		}

		if r.Sign() == 0 {
			return false
		}

		e := hashToInt(hash, c)
		s = new(big.Int).Mul(priv.D, r)
		s.Add(s, e)
		s.Mul(s, inv)
		s.Mod(s, N)

		// https://bitcoin.stackexchange.com/questions/68254/how-can-i-fix-this-non-canonical-signature-s-value-is-unnecessarily-high?rq=1
		// The signature is composed of two values, the r value and the s value.
		// If the s value is greater than N/2, which is not allowed.
		// Just add in some code that if s is greater than N/2, then s = N - s.
		if s.Cmp(halfOrder) == 1 {
			s.Sub(N, s)
		}
		return s.Sign() != 0
	})
	return
}

// copied from crypto/ecdsa.
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
