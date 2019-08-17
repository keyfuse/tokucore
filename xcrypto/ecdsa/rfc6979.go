// tokucore
//
// Copyright (c) 2014 Coda Hale
// Copyright 2019 by KeyFuse Labs
// BSD License

package ecdsa

import (
	"bytes"
	"hash"
	"math/big"

	"crypto/hmac"
	"crypto/sha256"
)

var (
	one = big.NewInt(1)
)

// NonceRFC6979 generates an ECDSA nonce (`k`) deterministically according to RFC 6979.
// https://tools.ietf.org/html/rfc6979#section-3.2
// It takes a 32-byte hash as an input and returns 32-byte nonce to be used in ECDSA algorithm.
func NonceRFC6979(q, x *big.Int, hash []byte) *big.Int {
	qlen := q.BitLen()
	alg := sha256.New
	holen := alg().Size()
	rolen := (qlen + 7) >> 3
	bx := append(int2octets(x, rolen), bits2octets(hash, q, qlen, rolen)...)

	// Step B
	v := bytes.Repeat([]byte{0x01}, holen)

	// Step C
	k := bytes.Repeat([]byte{0x00}, holen)

	// Step D
	k = mac(alg, k, append(append(v, 0x00), bx...), k)

	// Step E
	v = mac(alg, k, v, v)

	// Step F
	k = mac(alg, k, append(append(v, 0x01), bx...), k)

	// Step G
	v = mac(alg, k, v, v)

	// Step H
	for {
		// Step H1
		var t []byte

		// Step H2
		for len(t) < qlen/8 {
			v = mac(alg, k, v, v)
			t = append(t, v...)
		}

		// Step H3
		secret := bits2int(t, qlen)
		if secret.Cmp(one) >= 0 && secret.Cmp(q) < 0 {
			return secret
		}
		k = mac(alg, k, append(v, 0x00), k)
		v = mac(alg, k, v, v)
	}
}

// mac -- returns an HMAC of the given key and message
func mac(alg func() hash.Hash, k, m, buf []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(buf[:0])
}

// https://tools.ietf.org/html/rfc6979#section-2.3.2
func bits2int(in []byte, qlen int) *big.Int {
	vlen := len(in) * 8
	v := new(big.Int).SetBytes(in)
	if vlen > qlen {
		v = new(big.Int).Rsh(v, uint(vlen-qlen))
	}
	return v
}

// https://tools.ietf.org/html/rfc6979#section-2.3.3
func int2octets(v *big.Int, rolen int) []byte {
	out := v.Bytes()

	// pad with zeros if it's too short
	if len(out) < rolen {
		out2 := make([]byte, rolen)
		copy(out2[rolen-len(out):], out)
		return out2
	}

	// drop most significant bytes if it's too long
	if len(out) > rolen {
		out2 := make([]byte, rolen)
		copy(out2, out[len(out)-rolen:])
		return out2
	}

	return out
}

// https://tools.ietf.org/html/rfc6979#section-2.3.4
func bits2octets(in []byte, q *big.Int, qlen, rolen int) []byte {
	z1 := bits2int(in, qlen)
	z2 := new(big.Int).Sub(z1, q)
	if z2.Sign() < 0 {
		return int2octets(z1, rolen)
	}
	return int2octets(z2, rolen)
}
