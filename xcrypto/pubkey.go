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

// PublicKey -- an ecdsa.PublicKey with additional functions to
// serialize in uncompressed, compressed, and hybrid formats.
type PublicKey ecdsa.PublicKey

// PubKeyFromBytes -- parse bytes to public key.
func PubKeyFromBytes(key []byte) (*PublicKey, error) {
	curve := SECP256K1()
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
		x, y := SecUnmarshal(curve, key)
		if x == nil || y == nil {
			return nil, fmt.Errorf("pubkey.format.invalid:%v", format)
		}
		pubkey.X, pubkey.Y = x, y
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
	curve := p.Curve

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
	byteLen := (p.Curve.Params().BitSize + 7) >> 3

	// Format.
	format := pubkeyUncompressed
	key.WriteByte(format)

	xBytes := p.X.Bytes()
	xbuf := make([]byte, byteLen)
	copy(xbuf[byteLen-len(xBytes):], xBytes)
	key.Write(xbuf)

	yBytes := p.Y.Bytes()
	ybuf := make([]byte, byteLen)
	copy(ybuf[byteLen-len(yBytes):], yBytes)
	key.Write(ybuf)
	return key.Bytes()
}

// SerializeCompressed -- encoding a public key in a 33-byte compressed foramt.
func (p *PublicKey) SerializeCompressed() []byte {
	return SecMarshal(p.Curve, p.X, p.Y)
}

// Hash160 -- returns the Hash160 of the compressed public key.
func (p *PublicKey) Hash160() []byte {
	return Hash160(p.SerializeCompressed())
}

// SerializeHybrid -- encoding a public key in a 65-byte hybrid format.
func (p *PublicKey) SerializeHybrid() []byte {
	var key bytes.Buffer
	byteLen := (p.Curve.Params().BitSize + 7) >> 3

	// Format.
	format := pubkeyHybrid
	if isOdd(p.Y) {
		format |= 0x1
	}
	key.WriteByte(format)

	xBytes := p.X.Bytes()
	xbuf := make([]byte, byteLen)
	copy(xbuf[byteLen-len(xBytes):], xBytes)
	key.Write(xbuf)

	yBytes := p.Y.Bytes()
	ybuf := make([]byte, byteLen)
	copy(ybuf[byteLen-len(yBytes):], yBytes)
	key.Write(ybuf)
	return key.Bytes()
}
