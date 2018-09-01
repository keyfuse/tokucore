// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcrypto

import (
	"fmt"
	"math/big"

	"crypto/ecdsa"
	"encoding/asn1"
)

// Signature -- a type representing an ecdsa signature.
type Signature struct {
	R *big.Int
	S *big.Int
}

// Serialize -- encode Signature to the DER.
func (sig *Signature) Serialize() ([]byte, error) {
	der, err := asn1.Marshal(*sig)
	if err != nil {
		return nil, err
	}
	return der, nil
}

// Deserialize -- decoce the Signature from the DER encoding.
func (sig *Signature) Deserialize(der []byte) error {
	_, err := asn1.Unmarshal(der, sig)
	return err
}

// Verify --
// calls ecdsa.Verify to verify the signature of hash using the public key.
// It returns true if the signature is valid, false otherwise.
func (sig *Signature) Verify(hash []byte, pubKey *PublicKey) bool {
	return ecdsa.Verify(pubKey.ToECDSA(), hash, sig.R, sig.S)
}

// Sign -- sign a hash and return the signature with DER format.
func Sign(hash []byte, prv *PrivateKey) ([]byte, error) {
	sig, err := prv.Sign(hash)
	if err != nil {
		return nil, err
	}

	der, err := sig.Serialize()
	if err != nil {
		return nil, err
	}
	return der, nil
}

// Verify -- verify the signature with ECDSA public key.
func Verify(hash []byte, sign []byte, pub *PublicKey) error {
	sig := &Signature{}
	if err := sig.Deserialize(sign); err != nil {
		return err
	}
	if !sig.Verify(hash, pub) {
		return fmt.Errorf("signature.verify.failed")
	}
	return nil
}
