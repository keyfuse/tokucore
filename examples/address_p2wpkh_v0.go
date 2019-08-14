// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package main

import (
	"fmt"

	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xcore"
	"github.com/tokublock/tokucore/xcore/bip32"
)

func main() {
	seed := []byte("this.is.bohu.seed.")
	hdprvkey := bip32.NewHDKey(seed)
	pubkey := hdprvkey.PublicKey()

	addr := xcore.NewPayToWitnessV0PubKeyHashAddress(pubkey.Hash160())
	fmt.Printf("p2wpkh.address(mainet):\t%s\n", addr.ToString(network.MainNet))
}
