// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package main

import (
	"fmt"

	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xcore"
)

func main() {
	// A.
	seed := []byte("this.is.a.seed.")
	aHDKey := xcore.NewHDKey(seed)
	aPub := aHDKey.PublicKey().Serialize()

	// B.
	seed = []byte("this.is.b.seed.")
	bHDKey := xcore.NewHDKey(seed)
	bPub := bHDKey.PublicKey().Serialize()

	// C.
	seed = []byte("this.is.c.seed.")
	cHDKey := xcore.NewHDKey(seed)
	cPub := cHDKey.PublicKey().Serialize()

	// Redeem script.
	redeemScript := xcore.NewPayToMultiSigScript(2, aPub, bPub, cPub)
	multi := redeemScript.GetAddress()

	fmt.Printf("multisig.address(mainet):\t%s\n", multi.ToString(network.MainNet))
}
