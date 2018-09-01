// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcrypto

import (
	"bytes"
	"fmt"
	"math/big"

	"crypto/ecdsa"
)

const (
	// Length
	pubKeyBytesLenCompressed   = 33
	pubKeyBytesLenUncompressed = 65

	// Format
	pubkeyCompressed   byte = 0x2 // y_bit + x coord
	pubkeyUncompressed byte = 0x4 // x coord + y coord
	pubkeyHybrid       byte = 0x6 // y_bit + x coord + y coord
)

var (
	curve       = SECP256K1()
	curveParams = curve.Params()
)

// PublicKey -- an ecdsa.PublicKey with additional functions to
// serialize in uncompressed, compressed, and hybrid formats.
type PublicKey ecdsa.PublicKey

// PubKeyFromBytes -- parse bytes to public key.
func PubKeyFromBytes(key []byte) (*PublicKey, error) {
	pubkey := PublicKey{
		Curve: curve, // secp256k1 curve
	}

	pkLen := len(key)
	if pkLen == 0 {
		return nil, fmt.Errorf("pubkey.string.empty")
	}

	// Format.
	format := key[0]
	ybit := (format & 0x1) == 0x1
	format &= ^byte(0x1)

	switch pkLen {
	case pubKeyBytesLenUncompressed:
		// Format invalid.
		if format != pubkeyUncompressed && format != pubkeyHybrid {
			return nil, fmt.Errorf("pubkey.format.invalid:%v", format)
		}

		pubkey.X = new(big.Int).SetBytes(key[1:33])
		pubkey.Y = new(big.Int).SetBytes(key[33:])
		// ybit invalid.
		if format == pubkeyHybrid && ybit != isOdd(pubkey.Y) {
			return nil, fmt.Errorf("pubkey.ybit[%v].doesnt.match.oddness[%v]", ybit, isOdd(pubkey.Y))
		}
	case pubKeyBytesLenCompressed:
		if format != pubkeyCompressed {
			return nil, fmt.Errorf("pubkey.format.invalid:%v", format)
		}
		pubkey.X, pubkey.Y = expandPublicKey(key)
	default:
		return nil, fmt.Errorf("pubkey.size[%v].invalid", pkLen)
	}

	// Curve check.
	P := pubkey.Curve.Params().P
	if pubkey.X.Cmp(P) >= 0 {
		return nil, fmt.Errorf("pubkey.X[%v].is.greater.than.P[%v]", pubkey.X, P)
	}
	if pubkey.Y.Cmp(P) >= 0 {
		return nil, fmt.Errorf("pubkey.Y[%v].is.greater.than.P[%v]", pubkey.Y, P)
	}
	if !pubkey.Curve.IsOnCurve(pubkey.X, pubkey.Y) {
		return nil, fmt.Errorf("pubkey.point[%v, %v].is.not.on.the.curve", pubkey.X, pubkey.Y)
	}
	return &pubkey, nil
}

// ToECDSA -- returns the public key as a *ecdsa.PublicKey.
func (p *PublicKey) ToECDSA() *ecdsa.PublicKey {
	return (*ecdsa.PublicKey)(p)
}

// XBytes -- returns the x coord bytes.
func (p *PublicKey) XBytes() []byte {
	return p.X.Bytes()
}

// YBytes -- returns the y coord bytes.
func (p *PublicKey) YBytes() []byte {
	return p.Y.Bytes()
}

// Add -- add n2 to PublicKey.
func (p *PublicKey) Add(n2 []byte) *PublicKey {
	x1 := p.X
	y1 := p.Y

	// Private key.
	// pubkey1 = prvkey1*G
	// prvkey1 = n1
	// prvkey2 = n1 + n2
	// pubkey2 = prvkey2*G = (n1 + n2)*G
	privkey := PrvKeyFromBytes(n2)
	p2, err := PubKeyFromBytes(privkey.PubKey().SerializeUncompressed())
	if err != nil {
		panic(err)
	}
	x2 := p2.X
	y2 := p2.Y
	x3, y3 := curve.Add(x1, y1, x2, y2)
	return &PublicKey{
		X:     x3,
		Y:     y3,
		Curve: curve,
	}
}

// Serialize -- returns the compressed endcoding.
func (p *PublicKey) Serialize() []byte {
	return p.SerializeCompressed()
}

// SerializeUncompressed -- encoding public key in a 65-byte uncompressed format.
func (p *PublicKey) SerializeUncompressed() []byte {
	var key bytes.Buffer

	// Format.
	format := pubkeyUncompressed
	key.WriteByte(format)

	// X with padding to 32-bytes.
	xBytes := p.X.Bytes()
	for i := 0; i < (32 - len(xBytes)); i++ {
		key.WriteByte(0x0)
	}
	key.Write(xBytes)

	// Y with padding to 32-bytes.
	yBytes := p.Y.Bytes()
	for i := 0; i < (32 - len(yBytes)); i++ {
		key.WriteByte(0x0)
	}
	key.Write(yBytes)
	return key.Bytes()
}

// SerializeCompressed -- encoding a public key in a 33-byte compressed foramt.
func (p *PublicKey) SerializeCompressed() []byte {
	var key bytes.Buffer

	// Format.
	format := pubkeyCompressed
	if isOdd(p.Y) {
		format |= 0x1
	}
	key.WriteByte(format)

	// X with padding to 32-bytes.
	xBytes := p.X.Bytes()
	for i := 0; i < (32 - len(xBytes)); i++ {
		key.WriteByte(0x0)
	}
	key.Write(xBytes)
	return key.Bytes()
}

// Hash160 -- returns the Hash160 of the compressed public key.
func (p *PublicKey) Hash160() []byte {
	return Hash160(p.SerializeCompressed())
}

// SerializeHybrid -- encoding a public key in a 65-byte hybrid format.
func (p *PublicKey) SerializeHybrid() []byte {
	var key bytes.Buffer

	// Format.
	format := pubkeyHybrid
	if isOdd(p.Y) {
		format |= 0x1
	}
	key.WriteByte(format)

	// X with padding to 32-bytes.
	xBytes := p.X.Bytes()
	for i := 0; i < (32 - len(xBytes)); i++ {
		key.WriteByte(0x0)
	}
	key.Write(xBytes)

	// Y with padding to 32-bytes.
	yBytes := p.Y.Bytes()
	for i := 0; i < (32 - len(yBytes)); i++ {
		key.WriteByte(0x0)
	}
	key.Write(yBytes)
	return key.Bytes()
}

func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}

// As described at https://crypto.stackexchange.com/a/8916.
func expandPublicKey(key []byte) (*big.Int, *big.Int) {
	Y := big.NewInt(0)
	X := big.NewInt(0)
	X.SetBytes(key[1:])

	// y^2 = x^3 + ax^2 + b
	// a = 0
	// => y^2 = x^3 + b
	ySquared := big.NewInt(0)
	ySquared.Exp(X, big.NewInt(3), nil)
	ySquared.Add(ySquared, curveParams.B)

	Y.ModSqrt(ySquared, curveParams.P)

	Ymod2 := big.NewInt(0)
	Ymod2.Mod(Y, big.NewInt(2))

	signY := uint64(key[0]) - 2
	if signY != Ymod2.Uint64() {
		Y.Sub(curveParams.P, Y)
	}
	return X, Y
}
