// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package main

import (
	"fmt"

	"github.com/keyfuse/tokucore/xcore"
	"github.com/keyfuse/tokucore/xcore/bip32"
)

func assertNil(err error) {
	if err != nil {
		panic(err)
	}
}

// Demo for send coin to P2PKH address.
func main() {
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := bip32.NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohuPub := bohuHDKey.PublicKey()
	bohu := xcore.NewPayToPubKeyHashAddress(bohuPub.Hash160())

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := bip32.NewHDKey(seed)
	satoshiPub := satoshiHDKey.PublicKey()
	satoshi := xcore.NewPayToPubKeyHashAddress(satoshiPub.Hash160())

	// Output:
	bohuCoin := xcore.NewCoinBuilder().AddOutput(
		"faf1520f1d3e818fca695c2a903baa4a7eec4954f0b35aa01be1f2c1d2cfd802",
		0,
		1*xcore.Unit,
		"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
	).ToCoins()[0]

	tx, err := xcore.NewTransactionBuilder().
		AddCoin(bohuCoin).
		AddKeys(bohuPrv).
		To(satoshi, 666666).
		Then().
		SetChange(bohu).
		SendFees(10000).
		Then().
		Sign().
		BuildTransaction()
	assertNil(err)

	// Verify.
	err = tx.Verify()
	assertNil(err)

	fmt.Printf("%v\n", tx.ToString())
	fmt.Printf("p2pkh.txid:%s\n", tx.ID())
	signedTx := tx.Serialize()
	fmt.Printf("p2pkh.signed.tx:%x\n", signedTx)
}
