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
)

func main() {
	seed := []byte("this.is.bohu.seed.")
	hdprvkey := bip32.NewHDKey(seed)
	pubkey := hdprvkey.PublicKey()

	addr := xcore.NewPayToWitnessV0PubKeyHashAddress(pubkey.Hash160())
	fmt.Printf("p2wpkh.address(mainet):\t%s\n", addr.ToString(network.MainNet))
}
