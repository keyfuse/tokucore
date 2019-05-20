// tokucore
//
// Copyright (c) 2018-2019 TokuBlock
// BSD License

package schnorr

import (
	"errors"
	"math/big"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"

	"github.com/tokublock/tokucore/xcrypto/secp256k1"
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
func Sign(prv *ecdsa.PrivateKey, m [32]byte) ([64]byte, error) {
	sig := [64]byte{}
	curve := prv.Curve

	D := prv.D
	N := curve.Params().N

	// k' = int(hash(bytes(d) || m)) mod n
	d := intToByte(D)
	k0, err := getK0(m, d, N)
	if err != nil {
		return sig, err
	}

	// R = k'G
	Rx, Ry := curve.ScalarBaseMult(k0.Bytes())
	rX := intToByte(Rx)
	k := getK(curve, Ry, k0)

	// P = dG
	Px, Py := curve.ScalarBaseMult(d)

	// e = int(hash(bytes(x(R)) || bytes(dG) || m)) mod n
	e := getE(curve, m, Px, Py, rX)

	// ed
	ed := new(big.Int)
	ed.Mul(e, D)

	// s = k + ed
	s := new(big.Int)
	s.Add(k, ed)
	s.Mod(s, N)

	copy(sig[:32], rX)
	copy(sig[32:], intToByte(s))

	// Clean.
	zeroSlice(d)
	k.SetInt64(0)

	return sig, nil
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
func Verify(pub *ecdsa.PublicKey, m [32]byte, sig [64]byte) bool {
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

	e := getE(curve, m, pub.X, pub.Y, intToByte(r))
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

// PartyR -- returns R point.
func PartyR(prv *ecdsa.PrivateKey, m [32]byte) *Scalar {
	curve := prv.Curve
	N := curve.Params().N
	d := intToByte(prv.D)

	// k' = int(hash(bytes(d) || m)) mod n
	k0, err := getK0(m, d, N)
	if err != nil {
		return nil
	}

	// R = k'G
	rx, ry := curve.ScalarBaseMult(k0.Bytes())
	return NewScalar(curve, rx, ry)
}

// PartySign -- sign the m with aggregate pub and R.
func PartySign(prv *ecdsa.PrivateKey, m [32]byte, R *Scalar, pub *ecdsa.PublicKey) ([32]byte, error) {
	sig := [32]byte{}
	curve := prv.Curve

	D := prv.D
	N := curve.Params().N

	// k' = int(hash(bytes(d) || m)) mod n
	d := intToByte(D)
	k0, err := getK0(m, d, N)
	if err != nil {
		return sig, err
	}
	k := getK(curve, R.Y, k0)

	// e = int(hash(bytes(x(R)) || bytes(dG) || m)) mod n
	e := getE(curve, m, pub.X, pub.Y, intToByte(R.X))

	// ed
	ed := new(big.Int)
	ed.Mul(e, D)

	// s = k + ed
	s := new(big.Int)
	s.Add(k, ed)
	s.Mod(s, N)

	copy(sig[:], intToByte(s))

	// Clean.
	zeroSlice(d)
	k.SetInt64(0)

	return sig, nil
}

// PartyAggregate -- aggregate the signatures to one.
func PartyAggregate(curve elliptic.Curve, R *Scalar, sigs ...[32]byte) ([64]byte, error) {
	N := curve.Params().N
	sigFinal := [64]byte{}
	aggS := new(big.Int)

	for _, sig := range sigs {
		s := new(big.Int).SetBytes(sig[:])
		aggS.Add(aggS, s)
	}
	copy(sigFinal[:32], intToByte(R.X))
	copy(sigFinal[32:], intToByte(aggS.Mod(aggS, N)))
	return sigFinal, nil
}

func getK(curve elliptic.Curve, Ry, k0 *big.Int) *big.Int {
	P := curve.Params().P
	N := curve.Params().N

	if big.Jacobi(Ry, P) == 1 {
		return k0
	}
	return k0.Sub(N, k0)
}

func getK0(m [32]byte, d []byte, N *big.Int) (*big.Int, error) {
	hash := sha256.Sum256(append(d, m[:]...))
	i := new(big.Int).SetBytes(hash[:])
	k0 := i.Mod(i, N)
	if k0.Sign() == 0 {
		return nil, errors.New("k0 is zero")
	}
	return k0, nil
}

func getE(curve elliptic.Curve, m [32]byte, Px, Py *big.Int, rX []byte) *big.Int {
	N := curve.Params().N
	r := append(rX, secp256k1.SecMarshal(curve, Px, Py)...)
	r = append(r, m[:]...)
	h := sha256.Sum256(r)
	i := new(big.Int).SetBytes(h[:])
	return i.Mod(i, N)
}

func intToByte(i *big.Int) []byte {
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
