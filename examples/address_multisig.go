// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package main

import (
	"fmt"

	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xcore"
	"github.com/keyfuse/tokucore/xcore/bip32"
	"github.com/keyfuse/tokucore/xcrypto"
)

func main() {
	// A.
	seed := []byte("this.is.a.seed.")
	aHDKey := bip32.NewHDKey(seed)
	aPub := aHDKey.PublicKey().Serialize()

	// B.
	seed = []byte("this.is.b.seed.")
	bHDKey := bip32.NewHDKey(seed)
	bPub := bHDKey.PublicKey().Serialize()

	// C.
	seed = []byte("this.is.c.seed.")
	cHDKey := bip32.NewHDKey(seed)
	cPub := cHDKey.PublicKey().Serialize()

	// Redeem script.
	redeemScript := xcore.NewPayToMultiSigScript(2, aPub, bPub, cPub)
	redeem, _ := redeemScript.GetLockingScriptBytes()
	multi := xcore.NewPayToScriptHashAddress(xcrypto.Hash160(redeem))

	fmt.Printf("multisig.address(mainet):\t%s\n", multi.ToString(network.MainNet))
}
