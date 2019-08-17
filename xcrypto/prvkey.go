// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcrypto

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"

	"github.com/keyfuse/tokucore/xcrypto/secp256k1"
)

const (
	// PrvKeyBytesLen -- defines the length in bytes of a serialized private key.
	PrvKeyBytesLen = 32
)

// PrvKey --
type PrvKey ecdsa.PrivateKey

// PrvKeyFromBytes -- returns a private and public key for secp256k1 curve.
func PrvKeyFromBytes(key []byte) *PrvKey {
	curve := secp256k1.SECP256K1()
	x, y := curve.ScalarBaseMult(key)
	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(key),
	}
	return (*PrvKey)(priv)
}

// PubKey -- returns ecdsa public key.
func (p *PrvKey) PubKey() *PubKey {
	return (*PubKey)(&p.PublicKey)
}

// Add -- add n2 to PrvKey.
// k3 = (k1 + k2) mod N
func (p *PrvKey) Add(n2 []byte) *PrvKey {
	kint1 := new(big.Int).Set(p.D)
	kint2 := new(big.Int).SetBytes(n2)
	kint1.Add(kint1, kint2)
	kint1.Mod(kint1, p.Curve.Params().N)
	return PrvKeyFromBytes(kint1.Bytes())
}

// Serialize --
// returns the private key number d as a big-endian binary-encoded
// number, padded to a length of 32 bytes.
func (p *PrvKey) Serialize() []byte {
	var key bytes.Buffer

	dBytes := p.D.Bytes()
	for i := 0; i < (PrvKeyBytesLen - len(dBytes)); i++ {
		key.WriteByte(0x0)
	}
	key.Write(dBytes)
	return key.Bytes()
}
