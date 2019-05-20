// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcrypto

import (
	"encoding/asn1"
	"fmt"
	"math/big"

	"github.com/tokublock/tokucore/xcrypto/ecdsa"
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

// Sign -- sign a hash and return the signature with DER format.
func Sign(hash []byte, prv *PrivateKey) ([]byte, error) {
	r, s, err := ecdsa.Sign(prv.ToECDSA(), hash)
	if err != nil {
		return nil, err
	}
	sig := &Signature{R: r, S: s}
	return sig.Serialize()
}

// Verify -- verify the signature with ECDSA public key.
func Verify(hash []byte, sign []byte, pub *PublicKey) error {
	sig := &Signature{}
	if err := sig.Deserialize(sign); err != nil {
		return err
	}
	if !ecdsa.Verify(pub.ToECDSA(), hash, sig.R, sig.S) {
		return fmt.Errorf("signature.verify.failed")
	}
	return nil
}
