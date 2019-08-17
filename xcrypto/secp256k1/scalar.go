// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package secp256k1

import (
	"math/big"

	"crypto/elliptic"
)

// Scalar -- point.
type Scalar struct {
	X *big.Int `json:"X"`
	Y *big.Int `json:"Y"`
}

// NewScalar -- creates new Scalar.
func NewScalar(x *big.Int, y *big.Int) *Scalar {
	return &Scalar{
		X: x,
		Y: y,
	}
}

// Add -- add s2 to the s return new Scalar.
func (s *Scalar) Add(curve elliptic.Curve, s2 *Scalar) *Scalar {
	x3, y3 := curve.Add(s.X, s.Y, s2.X, s2.Y)
	return NewScalar(x3, y3)
}
