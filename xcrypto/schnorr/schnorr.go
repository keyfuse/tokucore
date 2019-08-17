// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package schnorr

import (
	"errors"
	"math/big"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"

	"github.com/keyfuse/tokucore/xcrypto/secp256k1"
)

// Sign -- signature with Schnorr, returning a 64 byte signature.
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
func Sign(prv *ecdsa.PrivateKey, m []byte) (*big.Int, *big.Int, error) {
	curve := prv.Curve

	D := prv.D
	N := curve.Params().N

	// k' = int(hash(bytes(d) || m)) mod n
	d := IntToByte(D)
	k0, err := GetK0(m, d, N)
	if err != nil {
		return nil, nil, err
	}

	// R = k'G
	Rx, Ry := curve.ScalarBaseMult(k0.Bytes())
	rX := IntToByte(Rx)
	k := GetK(curve, Ry, k0)

	// P = dG
	Px, Py := curve.ScalarBaseMult(d)

	// e = int(hash(bytes(x(R)) || bytes(dG) || m)) mod n
	e := GetE(curve, m, Px, Py, rX)

	// ed
	ed := new(big.Int)
	ed.Mul(e, D)

	// s = k + ed
	s := new(big.Int)
	s.Add(k, ed)
	s.Mod(s, N)

	// Clean.
	zeroSlice(d)
	k.SetInt64(0)

	return Rx, s, nil
}

// Verify -- verify the signature against the public key.
// https://github.com/sipa/bips/blob/bip-schnorr/bip-schnorr.mediawiki#verification
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
func Verify(pub *ecdsa.PublicKey, m []byte, r *big.Int, s *big.Int) bool {
	curve := pub.Curve
	P := curve.Params().P
	N := curve.Params().N

	if r.Cmp(P) > 0 || s.Cmp(N) > 0 {
		return false
	}

	e := GetE(curve, m, pub.X, pub.Y, IntToByte(r))
	sGx, sGy := curve.ScalarBaseMult(IntToByte(s))
	ePx, ePy := curve.ScalarMult(pub.X, pub.Y, IntToByte(e))

	// eP Inverse.
	ePy.Sub(P, ePy)

	// R=sG-eP=sG+(eP inverse)
	Rx, Ry := curve.Add(sGx, sGy, ePx, ePy)
	if (Rx.Sign() == 0 && Ry.Sign() == 0) || (big.Jacobi(Ry, P) != 1) || (Rx.Cmp(r) != 0) {
		return false
	}
	return true
}

// GetK0 -- used get k0 under schnorr BIP.
func GetK0(m []byte, d []byte, N *big.Int) (*big.Int, error) {
	hash := sha256.Sum256(append(d, m[:]...))
	i := new(big.Int).SetBytes(hash[:])
	k0 := i.Mod(i, N)
	if k0.Sign() == 0 {
		return nil, errors.New("k0 is zero")
	}
	return k0, nil
}

// GetK -- used get k under schnorr BIP.
func GetK(curve elliptic.Curve, Ry, k0 *big.Int) *big.Int {
	P := curve.Params().P
	N := curve.Params().N

	if big.Jacobi(Ry, P) == 1 {
		return k0
	}
	return k0.Sub(N, k0)
}

// GetE -- used get e under schnorr BIP.
func GetE(curve elliptic.Curve, m []byte, Px, Py *big.Int, rX []byte) *big.Int {
	N := curve.Params().N
	r := append(rX, secp256k1.SecMarshal(curve, Px, Py)...)
	r = append(r, m[:]...)
	h := sha256.Sum256(r)
	i := new(big.Int).SetBytes(h[:])
	return i.Mod(i, N)
}

// IntToByte -- used to convert the int to bytes under schnorr BIP.
func IntToByte(i *big.Int) []byte {
	b1, b2 := [32]byte{}, i.Bytes()
	copy(b1[32-len(b2):], b2)
	return b1[:]
}

// zeroSlice -- zeroes the memory of a scalar byte slice.
func zeroSlice(s []byte) {
	for i := 0; i < len(s); i++ {
		s[i] = 0x00
	}
}
