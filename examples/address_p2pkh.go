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
	hdpubkey := hdprvkey.HDPublicKey()
	addr := hdpubkey.GetAddress()

	fmt.Printf("p2pkh.address(mainet):\t%s\n", addr.ToString(network.MainNet))
	fmt.Printf("prv.wif(mainet):\t%s\n", hdprvkey.ToString(network.MainNet))
	fmt.Printf("pub.wif(mainet):\t%s\n", hdpubkey.ToString(network.MainNet))
}
