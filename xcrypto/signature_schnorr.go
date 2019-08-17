// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcrypto

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/keyfuse/tokucore/xcrypto/schnorr"
)

// SignatureSchnorr -- a type representing an schnorr signature.
type SignatureSchnorr struct {
	R *big.Int
	S *big.Int
}

// NewSignatureSchnorr -- create new SignatureSchnorr.
func NewSignatureSchnorr() *SignatureSchnorr {
	return &SignatureSchnorr{}
}

// Serialize -- used to serialize the struct to signature.
func (sig *SignatureSchnorr) Serialize() ([]byte, error) {
	sigFinal := make([]byte, 64)
	copy(sigFinal[:32], schnorr.IntToByte(sig.R))
	copy(sigFinal[32:], schnorr.IntToByte(sig.S))
	return sigFinal, nil
}

// Deserialize -- used to deserialize the signature to struct.
func (sig *SignatureSchnorr) Deserialize(sign []byte) error {
	sig.R = new(big.Int).SetBytes(sign[:32])
	sig.S = new(big.Int).SetBytes(sign[32:])
	return nil
}

// SchnorrSign -- used get the schnorr signature.
func SchnorrSign(prv *PrvKey, hash []byte) ([]byte, error) {
	eprv := (*ecdsa.PrivateKey)(prv)
	r, s, err := schnorr.Sign(eprv, hash)
	if err != nil {
		return nil, err
	}
	sig := &SignatureSchnorr{R: r, S: s}
	return sig.Serialize()
}

// SchnorrVerify -- used to verify the schnorr signature.
func SchnorrVerify(pub *PubKey, hash []byte, sign []byte) error {
	sig := NewSignatureSchnorr()
	if err := sig.Deserialize(sign); err != nil {
		return err
	}

	epub := (*ecdsa.PublicKey)(pub)
	if !schnorr.Verify(epub, hash, sig.R, sig.S) {
		return fmt.Errorf("schnorr.signature.verify.failed")
	}
	return nil
}
