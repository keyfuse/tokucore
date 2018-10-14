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
	seed := []byte("this.is.bohu.seed.")
	hdprvkey := xcore.NewHDKey(seed)
	pubkey := hdprvkey.PublicKey()

	addr := xcore.NewPayToWitnessPubKeyHashAddress(pubkey.Hash160())
	fmt.Printf("p2wpkh.address(mainet):\t%s\n", addr.ToString(network.MainNet))
}
