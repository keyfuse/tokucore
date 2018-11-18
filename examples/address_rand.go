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
	hdprvkey, _ := bip32.NewHDKeyRand()
	hdpubkey := hdprvkey.HDPublicKey()
	addr := xcore.NewPayToPubKeyHashAddress(hdprvkey.PublicKey().Hash160())

	fmt.Printf("address(mainet):\t%s\n", addr.ToString(network.MainNet))
	fmt.Printf("prv.wif(mainet):\t%s\n", hdprvkey.ToString(network.MainNet))
	fmt.Printf("pub.wif(mainet):\t%s\n", hdpubkey.ToString(network.MainNet))
}
