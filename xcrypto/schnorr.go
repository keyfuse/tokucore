// tokucore
//
// Copyright (c) 2018-2019 TokuBlock
// BSD License

package xcrypto

import (
	"errors"
	"math/big"

	"crypto/ecdsa"
	"crypto/elliptic"
)

// SchnorrSign -- signature with Schnorr, returning a 64 byte signature..
// https://github.com/sipa/bips/blob/bip-schnorr/bip-schnorr.mediawiki#signing
func SchnorrSign(priv *ecdsa.PrivateKey, hash []byte) ([64]byte, error) {
	sig := [64]byte{}
	D := priv.D
	curve := priv.Curve
	P := curve.Params().P
	N := curve.Params().N

	// k' = int(hash(bytes(d) || m)) mod n
	d := intToByte(D)
	k0, err := getK0(hash, d, N)
	if err != nil {
		return sig, err
	}

	// R = k'G
	Rx, Ry := curve.ScalarBaseMult(k0.Bytes())
	k := getK(Ry, k0, P, N)

	// dG
	Px, Py := curve.ScalarBaseMult(d)

	// bytes(Rx)
	rX := intToByte(Rx)

	// e = int(hash(bytes(x(R)) || bytes(dG) || m)) mod n
	e := getE(curve, Px, Py, rX, hash, N)
	e.Mul(e, D)
	k.Add(k, e)
	k.Mod(k, N)

	copy(sig[:32], rX)
	copy(sig[32:], intToByte(k))
	return sig, nil
}

func intToByte(i *big.Int) []byte {
	b1, b2 := [32]byte{}, i.Bytes()
	copy(b1[32-len(b2):], b2)
	return b1[:]
}

func getK(Ry, k0 *big.Int, P *big.Int, N *big.Int) *big.Int {
	if big.Jacobi(Ry, P) == 1 {
		return k0
	}
	return k0.Sub(N, k0)
}

func getK0(hash []byte, d []byte, N *big.Int) (*big.Int, error) {
	h := Sha256(append(d, hash[:]...))
	i := new(big.Int).SetBytes(h[:])
	k0 := i.Mod(i, N)
	if k0.Sign() == 0 {
		return nil, errors.New("k0 is zero")
	}
	return k0, nil
}

func getE(curve elliptic.Curve, Px, Py *big.Int, rX []byte, m []byte, N *big.Int) *big.Int {
	r := append(rX, SecMarshal(curve, Px, Py)...)
	r = append(r, m[:]...)
	h := Sha256(r)
	i := new(big.Int).SetBytes(h[:])
	return i.Mod(i, N)
}
