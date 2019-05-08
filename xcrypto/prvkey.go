// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcrypto

import (
	"bytes"
	"math/big"

	"crypto/ecdsa"
)

const (
	// PrvKeyBytesLen -- defines the length in bytes of a serialized private key.
	PrvKeyBytesLen = 32
)

// PrivateKey --
type PrivateKey ecdsa.PrivateKey

// PrvKeyFromBytes -- returns a private and public key for secp256k1 curve.
func PrvKeyFromBytes(key []byte) *PrivateKey {
	curve := SECP256K1()
	x, y := curve.ScalarBaseMult(key)
	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(key),
	}
	return (*PrivateKey)(priv)
}

// PubKey -- returns ecdsa public key.
func (p *PrivateKey) PubKey() *PublicKey {
	return (*PublicKey)(&p.PublicKey)
}

// Add -- add n2 to PrivateKey.
// k3 = (k1 + k2) mod N
func (p *PrivateKey) Add(n2 []byte) *PrivateKey {
	kint1 := p.D
	kint2 := new(big.Int).SetBytes(n2)
	kint1.Add(kint1, kint2)
	kint1.Mod(kint1, p.Curve.Params().N)
	return PrvKeyFromBytes(kint1.Bytes())
}

// ToECDSA -- returns the private key as a *ecdsa.PrivateKey.
func (p *PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	return (*ecdsa.PrivateKey)(p)
}

// Serialize --
// returns the private key number d as a big-endian binary-encoded
// number, padded to a length of 32 bytes.
func (p *PrivateKey) Serialize() []byte {
	var key bytes.Buffer

	dBytes := p.ToECDSA().D.Bytes()
	for i := 0; i < (PrvKeyBytesLen - len(dBytes)); i++ {
		key.WriteByte(0x0)
	}
	key.Write(dBytes)
	return key.Bytes()
}
