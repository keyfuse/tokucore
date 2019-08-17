// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package main

import (
	"fmt"

	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xcore/bip32"
	"github.com/keyfuse/tokucore/xcore/bip39"
)

func main() {
	net := network.TestNet
	mnemonic, err := bip39.NewBIP39(bip39.CHINESE)
	if err != nil {
		panic(err)
	}

	hdkey := bip32.NewHDKey(mnemonic.Seed())
	fmt.Printf("mnemonic:\t%v\n", mnemonic.ToString())
	fmt.Printf("bip39_seed:\t%x\n", mnemonic.Seed())
	fmt.Printf("bip32_xprv:\t%v\n", hdkey.ToString(net))
}
