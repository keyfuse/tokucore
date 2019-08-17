// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcrypto

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/keyfuse/tokucore/xcrypto/secp256k1"
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

// PubKey -- an ecdsa.PubKey with additional functions to
// serialize in uncompressed, compressed, and hybrid formats.
type PubKey ecdsa.PublicKey

// PubKeyFromBytes -- parse bytes to public key.
func PubKeyFromBytes(key []byte) (*PubKey, error) {
	curve := secp256k1.SECP256K1()
	pubkey := PubKey{
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
		isOdd := (pubkey.Y.Bit(0) == 1)
		// ybit invalid.
		if format == pubkeyHybrid && ybit != isOdd {
			return nil, fmt.Errorf("pubkey.ybit[%v].doesnt.match.oddness[%v]", ybit, pubkey.Y.Bit(0))
		}
	case pubKeyBytesLenCompressed:
		x, y := secp256k1.SecUnmarshal(curve, key)
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

// XBytes -- returns the x coord bytes.
func (p *PubKey) XBytes() []byte {
	return p.X.Bytes()
}

// YBytes -- returns the y coord bytes.
func (p *PubKey) YBytes() []byte {
	return p.Y.Bytes()
}

// Add -- add p2 to PubKey.
func (p *PubKey) Add(p2 *PubKey) *PubKey {
	x1 := p.X
	y1 := p.Y
	curve := p.Curve

	x2 := p2.X
	y2 := p2.Y
	x3, y3 := curve.Add(x1, y1, x2, y2)
	return &PubKey{
		X:     x3,
		Y:     y3,
		Curve: curve,
	}
}

// Serialize -- returns the compressed endcoding.
func (p *PubKey) Serialize() []byte {
	return p.SerializeCompressed()
}

// SerializeUncompressed -- encoding public key in a 65-byte uncompressed format.
func (p *PubKey) SerializeUncompressed() []byte {
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
func (p *PubKey) SerializeCompressed() []byte {
	return secp256k1.SecMarshal(p.Curve, p.X, p.Y)
}

// Hash160 -- returns the Hash160 of the compressed public key.
func (p *PubKey) Hash160() []byte {
	return Hash160(p.SerializeCompressed())
}

// SerializeHybrid -- encoding a public key in a 65-byte hybrid format.
func (p *PubKey) SerializeHybrid() []byte {
	var key bytes.Buffer
	byteLen := (p.Curve.Params().BitSize + 7) >> 3

	// Format.
	format := pubkeyHybrid
	if p.Y.Bit(0) == 1 {
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
