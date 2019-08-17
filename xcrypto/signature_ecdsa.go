// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcrypto

import (
	"crypto/ecdsa"
	"encoding/asn1"
	"fmt"
	"math/big"

	xecdsa "github.com/keyfuse/tokucore/xcrypto/ecdsa"
)

// SignatureEcdsa -- a type representing an ECDSA signature.
type SignatureEcdsa struct {
	R *big.Int
	S *big.Int
}

// NewSignatureEcdsa -- create new SignatureEcdsa.
func NewSignatureEcdsa() *SignatureEcdsa {
	return &SignatureEcdsa{}
}

// Serialize -- used to serialize the struct to signature.
func (sig *SignatureEcdsa) Serialize() ([]byte, error) {
	der, err := asn1.Marshal(*sig)
	if err != nil {
		return nil, err
	}
	return der, nil
}

// Deserialize -- used to deserialize the signature to struct.
func (sig *SignatureEcdsa) Deserialize(sign []byte) error {
	_, err := asn1.Unmarshal(sign, sig)
	return err
}

// EcdsaSign -- used get the ecdsa signature.
func EcdsaSign(prv *PrvKey, hash []byte) ([]byte, error) {
	eprv := (*ecdsa.PrivateKey)(prv)
	r, s, err := xecdsa.Sign(eprv, hash)
	if err != nil {
		return nil, err
	}
	sig := &SignatureEcdsa{R: r, S: s}
	return sig.Serialize()
}

// EcdsaVerify -- used to verify the ecdsa signature.
func EcdsaVerify(pub *PubKey, hash []byte, sign []byte) error {
	sig := NewSignatureEcdsa()
	if err := sig.Deserialize(sign); err != nil {
		return err
	}

	epub := (*ecdsa.PublicKey)(pub)
	if !xecdsa.Verify(epub, hash, sig.R, sig.S) {
		return fmt.Errorf("ecdsa.signature.verify.failed")
	}
	return nil
}
