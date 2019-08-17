// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package secp256k1

import (
	"crypto/elliptic"
	"math/big"
)

// SecMarshal -- convert a point into the form specified in section 2.3.3 of the SEC 1 standard.
// https://www.ipa.go.jp/security/enc/CRYPTREC/fy15/doc/1_01sec1.pdf
func SecMarshal(curve elliptic.Curve, x *big.Int, y *big.Int) []byte {
	byteLen := (curve.Params().BitSize + 7) >> 3
	ret := make([]byte, 1+byteLen)

	// 0x02 for even, 0x03 for odd.
	format := byte(0x02)
	// is Odd.
	if y.Bit(0) == 1 {
		format |= 0x01
	}
	ret[0] = format
	xBytes := x.Bytes()
	copy(ret[1+byteLen-len(xBytes):], xBytes)
	return ret
}

// SecUnmarshal -- converts a point, serialised by Marshal, into an x, y pair.
// On error, x = nil.
// As described at https://crypto.stackexchange.com/a/8916.
func SecUnmarshal(curve elliptic.Curve, data []byte) (*big.Int, *big.Int) {
	Y := big.NewInt(0)
	X := big.NewInt(0)
	curveParams := curve.Params()

	// Check format.
	if (data[0] &^ 1) != 2 {
		return nil, nil
	}

	// y^2 = x^3 + b
	X.SetBytes(data[1:])
	ySquared := big.NewInt(0)
	ySquared.Exp(X, three, nil)
	ySquared.Add(ySquared, curveParams.B)
	Y.ModSqrt(ySquared, curveParams.P)

	// Y is odd or even.
	if Y.Bit(0) != uint(data[0]&1) {
		Y.Sub(curveParams.P, Y)
	}
	return X, Y
}
