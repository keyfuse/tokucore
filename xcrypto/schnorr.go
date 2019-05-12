// tokucore
//
// Copyright (c) 2018-2019 TokuBlock
// BSD License

package xcrypto

import (
	"errors"
	"math/big"

	"crypto/elliptic"
)

// SchnorrSign -- signature with Schnorr, returning a 64 byte signature.
// https://github.com/sipa/bips/blob/bip-schnorr/bip-schnorr.mediawiki#signing
// Input:
//   The secret key d: an integer in the range [1..n-1].
//   The message m: a 32-byte array
//
// To sign m for public key dG:
//   Let k' = int(hash(bytes(d) || m)) mod n
//   Fail if k' = 0
//   Let R = k'G
//   Let k = k' if jacobi(y(R)) = 1, otherwise let k = n - k'
//   Let e = int(hash(bytes(x(R)) || bytes(dG) || m)) mod n
//   The signature is bytes(x(R)) || bytes((k + ed) mod n)
func SchnorrSign(priv *PrivateKey, m [32]byte) ([64]byte, error) {
	sig := [64]byte{}
	curve := priv.Curve
	D := priv.D
	P := curve.Params().P
	N := curve.Params().N

	// k' = int(hash(bytes(d) || m)) mod n
	d := intToByte(D)
	k0, err := getK0(m, d, N)
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
	e := getE(curve, m, Px, Py, rX, N)
	e.Mul(e, D)
	k.Add(k, e)
	k.Mod(k, N)

	copy(sig[:32], rX)
	copy(sig[32:], intToByte(k))
	return sig, nil
}

// https://github.com/sipa/bips/blob/bip-schnorr/bip-schnorr.mediawiki#signing
// Input:
//   The public key pk: a 33-byte array
//   The message m: a 32-byte array
//
// A signature sig: a 64-byte array
//   The signature is valid if and only if the algorithm below does not fail
//   Let P = point(pk); fail if point(pk) fails
//   Let r = int(sig[0:32]); fail if r ≥ p
//   Let s = int(sig[32:64]); fail if s ≥ n
//   Let e = int(hash(bytes(r) || bytes(P) || m)) mod n
//   Let R = sG - eP
//   Fail if infinite(R)
//   Fail if jacobi(y(R)) ≠ 1 or x(R) ≠ r
func SchnorrVerify(pub *PublicKey, m [32]byte, sig [64]byte) bool {
	curve := pub.Curve
	P := curve.Params().P
	N := curve.Params().N

	r := new(big.Int).SetBytes(sig[:32])
	if r.Cmp(P) > 0 {
		return false
	}

	s := new(big.Int).SetBytes(sig[32:])
	if s.Cmp(N) > 0 {
		return false
	}

	e := getE(curve, m, pub.X, pub.Y, intToByte(r), N)
	sGx, sGy := curve.ScalarBaseMult(intToByte(s))
	ePx, ePy := curve.ScalarMult(pub.X, pub.Y, intToByte(e))

	// eP Inverse.
	ePy.Sub(P, ePy)

	// R= sG - eP= sG + (eP inverse)
	Rx, Ry := curve.Add(sGx, sGy, ePx, ePy)
	if (Rx.Sign() == 0 && Ry.Sign() == 0) || (big.Jacobi(Ry, P) != 1) || (Rx.Cmp(r) != 0) {
		return false
	}
	return true
}

func getK(Ry, k0 *big.Int, P *big.Int, N *big.Int) *big.Int {
	if big.Jacobi(Ry, P) == 1 {
		return k0
	}
	return k0.Sub(N, k0)
}

func getK0(m [32]byte, d []byte, N *big.Int) (*big.Int, error) {
	hash := Sha256(append(d, m[:]...))
	i := new(big.Int).SetBytes(hash[:])
	k0 := i.Mod(i, N)
	if k0.Sign() == 0 {
		return nil, errors.New("k0 is zero")
	}
	return k0, nil
}

func getE(curve elliptic.Curve, m [32]byte, Px, Py *big.Int, rX []byte, N *big.Int) *big.Int {
	r := append(rX, SecMarshal(curve, Px, Py)...)
	r = append(r, m[:]...)
	h := Sha256(r)
	i := new(big.Int).SetBytes(h[:])
	return i.Mod(i, N)
}

func intToByte(i *big.Int) []byte {
	b1, b2 := [32]byte{}, i.Bytes()
	copy(b1[32-len(b2):], b2)
	return b1[:]
}
