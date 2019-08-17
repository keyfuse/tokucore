// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package main

import (
	"fmt"

	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xcore/bip32"
)

// Bitcoin HD wallet demo.
func main() {
	net := network.TestNet
	seed := []byte("bitcoin blockchain tokublock sandbox")
	hdkey := bip32.NewHDKey(seed)

	// Master Private Key.
	masterprv := hdkey.ToString(net)
	fmt.Printf("master.prvkey:%v\n", masterprv)

	// bitcoin  path: m/44'/0'/0'/0
	{
		for i := 0; i < 2; i++ {
			path := fmt.Sprintf("m/44'/0'/0'/%d", i)
			prvkey, err := hdkey.DeriveByPath(path)
			if err != nil {
				panic(err)
			}
			fmt.Printf("btc.chain:%v\n", path)
			fmt.Printf("\tchild.prvkey:%v\n", prvkey.ToString(net))

			pubkey := prvkey.HDPublicKey()
			fmt.Printf("\tchild.pubkey:%v\n", pubkey.ToString(net))
		}
	}

	// ethereum path: m/44'/60'/0'/0
	{
		for i := 0; i < 2; i++ {
			path := fmt.Sprintf("m/44'/60'/0'/%d", i)
			prvkey, err := hdkey.DeriveByPath(path)
			if err != nil {
				panic(err)
			}
			fmt.Printf("eth.chain:%v\n", path)
			fmt.Printf("\tchild.prvkey:%v\n", prvkey.ToString(net))

			pubkey := prvkey.HDPublicKey()
			fmt.Printf("\tchild.pubkey:%v\n", pubkey.ToString(net))
		}
	}
}
