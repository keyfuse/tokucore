// tokucore
//
// Copyright (c) 2019 TokuBlock
// BSD License

package main

import (
	"fmt"
	"math/big"

	"github.com/tokublock/tokucore/xcrypto"
)

func assertNil(err error) {
	if err != nil {
		panic(err)
	}
}

// Demo for scriptless ECDSA adaptor signature.
func main() {
	hash := xcrypto.DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	// Party 1.
	p1, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	prv1 := xcrypto.PrvKeyFromBytes(p1.Bytes())
	pub1 := prv1.PubKey()
	alice, err := xcrypto.NewEcdsaAlice(prv1)
	assertNil(err)
	encpk1 := alice.EncPk()
	encpub1 := alice.EncPub()

	// Party 2.
	secret := new(big.Int).SetInt64(2019)
	p2, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c9", 16)
	prv2 := xcrypto.PrvKeyFromBytes(p2.Bytes())
	pub2 := prv2.PubKey()
	bob, err := xcrypto.NewEcdsaBob(prv2, secret)
	assertNil(err)
	encpk2 := bob.EncPk()
	encpub2 := bob.EncPub()

	// Phase 1.
	sharepub1 := alice.ScriptlessPhase1(pub2)
	sharepub2 := bob.ScriptlessPhase1(pub1)

	// Phase 2.
	scalarR1 := alice.ScriptlessPhase2(hash)
	scalarR2 := bob.ScriptlessPhase2(hash)

	// Phase 3.
	shareR1 := alice.ScriptlessPhase3(encpk2, encpub2, scalarR2)
	shareR2 := bob.ScriptlessPhase3(encpk1, encpub1, scalarR1)

	// Phase 4.
	sig1, err := alice.ScriptlessPhase4(shareR1)
	assertNil(err)
	sig2, err := bob.ScriptlessPhase4(shareR2)
	assertNil(err)

	// Phase 5.
	fs1, err := alice.ScriptlessPhase5(shareR1, sig2)
	assertNil(err)
	fs2, err := bob.ScriptlessPhase5(shareR2, sig1)
	assertNil(err)

	// Alice Phase 6.
	ft := alice.ScriptlessPhase6(fs1, fs2)

	// Bob Phase 6.
	dersig, err := bob.ScriptlessPhase6(shareR2, fs2)

	// Verify.
	err = xcrypto.EcdsaVerify(sharepub2, hash, dersig)
	assertNil(err)

	fmt.Printf("\nAdaptor secret: %x\nKeys\n  x1: %x\n  x2: %x\n  Q:  %x\n\nSignatures\n  %x\nIs valid under Q?: %v\nRecovered secret: %x\n",
		secret.Bytes(),
		p1.Bytes(),
		p2.Bytes(),
		sharepub1.SerializeCompressed(),
		fs1,
		err == nil,
		ft.Bytes())
}
