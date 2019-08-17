// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcrypto

import (
	"math/big"

	"crypto/elliptic"

	"github.com/keyfuse/tokucore/xcrypto/schnorr"
	"github.com/keyfuse/tokucore/xcrypto/secp256k1"
)

// SchnorrParty -- Schnorr party struct.
type SchnorrParty struct {
	k0    *big.Int
	N     *big.Int
	prv   *PrvKey
	pub   *PubKey
	hash  []byte
	curve elliptic.Curve
	r     *secp256k1.Scalar
}

// NewSchnorrParty -- creates new SchnorrParty.
func NewSchnorrParty(prv *PrvKey) (*SchnorrParty, error) {
	pub := prv.PubKey()
	curve := pub.Curve
	N := curve.Params().N
	return &SchnorrParty{
		N:     N,
		prv:   prv,
		pub:   pub,
		curve: curve,
	}, nil
}

// Phase1 -- used to generate final pubkey of parties.
// Return the shared PubKey.
func (party *SchnorrParty) Phase1(pub2 *PubKey) *PubKey {
	pub := party.pub
	return pub.Add(pub2)
}

// Phase2 -- used to generate k, kinv, scalarR.
// Return the party scalar R.
func (party *SchnorrParty) Phase2(hash []byte) *secp256k1.Scalar {
	N := party.N
	prv := party.prv
	pub := prv.PubKey()
	curve := pub.Curve
	d := schnorr.IntToByte(prv.D)

	party.hash = hash
	// Scalar R.
	// k' = int(hash(bytes(d) || m)) mod n
	k0, err := schnorr.GetK0(hash, d, N)
	if err != nil {
		return nil
	}
	party.k0 = k0

	rx, ry := curve.ScalarBaseMult(k0.Bytes())
	party.r = secp256k1.NewScalar(rx, ry)
	return party.r
}

// Phase3 -- return shared scalar R.
func (party *SchnorrParty) Phase3(r2 *secp256k1.Scalar) *secp256k1.Scalar {
	curve := party.curve
	scalarR := party.r

	shareScalarR := secp256k1.NewScalar(scalarR.X, scalarR.Y)
	return shareScalarR.Add(curve, r2)
}

// Phase4 -- return the signature of this party.
func (party *SchnorrParty) Phase4(sharePub *PubKey, shareR *secp256k1.Scalar) ([]byte, error) {
	k0 := party.k0
	m := party.hash
	N := party.N
	prv := party.prv
	pub := sharePub
	curve := party.curve
	scalarR := party.r
	shareScalarR := shareR

	// e = int(hash(bytes(x(R)) || bytes(dG) || m)) mod n
	e := schnorr.GetE(curve, m, pub.X, pub.Y, schnorr.IntToByte(scalarR.X))

	// ed
	ed := new(big.Int)
	ed.Mul(e, prv.D)

	// s = k + ed
	k := schnorr.GetK(curve, shareScalarR.Y, k0)
	s := new(big.Int)
	s.Add(k, ed)
	s.Mod(s, N)

	return schnorr.IntToByte(s), nil
}

// Phase5 -- return the final signature.
func (party *SchnorrParty) Phase5(shareR *secp256k1.Scalar, sigs ...[]byte) ([]byte, error) {
	N := party.N
	R := shareR

	aggs := new(big.Int)
	sigFinal := make([]byte, 64)

	for _, sig := range sigs {
		s := new(big.Int).SetBytes(sig[:])
		aggs.Add(aggs, s)
	}
	aggs = aggs.Mod(aggs, N)

	copy(sigFinal[:32], schnorr.IntToByte(R.X))
	copy(sigFinal[32:], schnorr.IntToByte(aggs))
	return sigFinal, nil
}

// Close -- close the party.
func (party *SchnorrParty) Close() {
	party.prv = nil
	if party.k0 != nil {
		party.k0.SetInt64(0)
	}
}
