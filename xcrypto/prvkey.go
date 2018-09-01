// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcrypto

import (
	"bytes"
	"math/big"

	"crypto/ecdsa"
	"crypto/sha256"
)

const (
	// PrvKeyBytesLen -- defines the length in bytes of a serialized private key.
	PrvKeyBytesLen = 32
)

// PrivateKey --
type PrivateKey ecdsa.PrivateKey

// PrvKeyFromBytes -- returns a private and public key for secp256k1 curve.
func PrvKeyFromBytes(key []byte) *PrivateKey {
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
	kint1.Mod(kint1, curveParams.N)
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

// Sign --
// generates an ECDSA signature for the provided hash (which should be the result
// of hashing a larger message) using the private key. Produced signature
// is deterministic (same message and same key yield the same signature) and canonical
// in accordance with RFC6979 and BIP0062.
func (p *PrivateKey) Sign(hash []byte) (*Signature, error) {
	sig := &Signature{}
	r, s, err := EcdsaSign(p.ToECDSA(), hash, sha256.New)
	if err != nil {
		return nil, err
	}
	sig.R = r
	sig.S = s
	return sig, nil
}
