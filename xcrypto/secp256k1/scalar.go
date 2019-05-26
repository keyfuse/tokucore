// tokucore
//
// Copyright (c) 2019 TokuBlock
// BSD License

package secp256k1

import (
	"math/big"

	"crypto/elliptic"
)

// Scalar -- point.
type Scalar struct {
	X     *big.Int
	Y     *big.Int
	curve elliptic.Curve
}

// NewScalar -- creates new Scalar.
func NewScalar(curve elliptic.Curve, x *big.Int, y *big.Int) *Scalar {
	return &Scalar{
		X:     x,
		Y:     y,
		curve: curve,
	}
}

// Add -- add s2 to the s return new Scalar.
func (s *Scalar) Add(s2 *Scalar) *Scalar {
	curve := s.curve
	x3, y3 := curve.Add(s.X, s.Y, s2.X, s2.Y)
	return NewScalar(curve, x3, y3)
}
