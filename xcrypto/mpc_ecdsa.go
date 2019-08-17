// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcrypto

import (
	"errors"
	"math/big"

	"crypto/elliptic"

	xecdsa "github.com/keyfuse/tokucore/xcrypto/ecdsa"
	"github.com/keyfuse/tokucore/xcrypto/paillier"
	"github.com/keyfuse/tokucore/xcrypto/secp256k1"
)

const (
	bitlen = 2048
)

// EcdsaParty -- ECDSA party struct.
type EcdsaParty struct {
	k      *big.Int
	kinv   *big.Int
	N      *big.Int
	prv    *PrvKey
	pub    *PubKey
	hash   []byte
	curve  elliptic.Curve
	encpk  *big.Int
	encprv *paillier.PrvKey
}

// NewEcdsaParty -- creates new EcdsaParty.
func NewEcdsaParty(prv *PrvKey) *EcdsaParty {
	pub := prv.PubKey()
	curve := pub.Curve
	N := curve.Params().N

	return &EcdsaParty{
		N:     N,
		prv:   prv,
		pub:   pub,
		curve: curve,
	}
}

// Phase1 -- used to generate final pubkey of parties.
// Return the shared PubKey.
func (party *EcdsaParty) Phase1(pub2 *PubKey) *PubKey {
	prv := party.prv
	pub := prv.PubKey()
	curve := pub.Curve

	px, py := curve.ScalarMult(pub2.X, pub2.Y, prv.D.Bytes())
	return &PubKey{X: px, Y: py, Curve: curve}
}

// Phase2 -- used to generate k, kinv, scalarR.
// Return the party scalar R.
func (party *EcdsaParty) Phase2(hash []byte) (*big.Int, *paillier.PubKey, *secp256k1.Scalar) {
	N := party.N
	prv := party.prv
	pub := prv.PubKey()
	curve := pub.Curve
	party.hash = hash

	// Paillier key pair.
	encpub, encprv, err := paillier.GenerateKeyPair(bitlen)
	if err != nil {
		return nil, nil, nil
	}
	party.encprv = encprv

	// Homomorphic Encryption of party pk.
	encpk, err := encpub.Encrypt(prv.D)
	if err != nil {
		return nil, nil, nil
	}
	party.encpk = encpk

	// RFC6979 K nonce.
	k := xecdsa.NonceRFC6979(N, prv.D, hash)
	kinv := new(big.Int).ModInverse(k, N)
	party.k = k
	party.kinv = kinv

	rx, ry := curve.ScalarBaseMult(k.Bytes())
	return encpk, encpub, secp256k1.NewScalar(rx, ry)
}

// Phase3 -- set party2's r2 to this party.
// Return the shared R.
func (party *EcdsaParty) Phase3(r2 *secp256k1.Scalar) *secp256k1.Scalar {
	k := party.k
	curve := party.curve
	rx, ry := curve.ScalarMult(r2.X, r2.Y, k.Bytes())

	return secp256k1.NewScalar(rx, ry)
}

// Phase4 -- generate the homomorphic encryption signature of this party.
// Return the homomorphic ciphertext.
func (party *EcdsaParty) Phase4(encpk2 *big.Int, encpub2 *paillier.PubKey, shareR *secp256k1.Scalar) (*big.Int, error) {
	var err error
	var ct *big.Int

	prv := party.prv
	pk1 := prv.D
	kinv := party.kinv
	hash := party.hash
	curve := party.curve

	// s’=(z+r⋅e(pk2)⋅pk1)/k1
	z := xecdsa.HashToInt(curve, hash)

	// e(pk2)⋅pk1
	if ct, err = encpub2.MultPlaintext(encpk2, pk1); err != nil {
		return nil, err
	}

	// r⋅e(pk2)⋅pk1
	if ct, err = encpub2.MultPlaintext(ct, shareR.X); err != nil {
		return nil, err
	}

	// z+r⋅e(pk2)⋅pk1
	if ct, err = encpub2.AddPlaintext(ct, z); err != nil {
		return nil, err
	}

	// (z+r⋅e(pk2)⋅pk1)/k1
	if ct, err = encpub2.MultPlaintext(ct, kinv); err != nil {
		return nil, err
	}
	return ct, nil
}

// Phase5 -- generate the final signature of two party.
// Return the final signature.
func (party *EcdsaParty) Phase5(shareR *secp256k1.Scalar, sign2 *big.Int) ([]byte, error) {
	N := party.N
	kinv := party.kinv
	encprv := party.encprv

	sig, err := encprv.Decrypt(sign2)
	if err != nil {
		return nil, err
	}
	s := sig.Mul(sig, kinv).Mod(sig, N)
	halfOrder := new(big.Int).Rsh(N, 1)
	if s.Cmp(halfOrder) == 1 {
		s.Sub(N, s)
	}
	if s.Sign() == 0 {
		return nil, errors.New("calculated S is zero")
	}
	esig := NewSignatureEcdsa()
	esig.R = shareR.X
	esig.S = s
	return esig.Serialize()
}

// Close -- used to cleanup the secret.
func (party *EcdsaParty) Close() {
	party.prv = nil
	party.encprv = nil
	if party.k != nil {
		party.k.SetInt64(0)
		party.kinv.SetInt64(0)
	}
}
