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

func assertNil(err error) {
	if err != nil {
		panic(err)
	}
}

// Demo for sent coin to Native SegWit P2WSH address and spending from SegWit address to normal address.
func main() {
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := bip32.NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohuPub := bohuHDKey.PublicKey()
	bohu := xcore.NewPayToPubKeyHashAddress(bohuPub.Hash160())

	// A.
	seed = []byte("this.is.a.seed.")
	aHDKey := bip32.NewHDKey(seed)
	aPrv := aHDKey.PrivateKey()
	aPub := aHDKey.PublicKey().Serialize()

	// B.
	seed = []byte("this.is.b.seed.")
	bHDKey := bip32.NewHDKey(seed)
	bPub := bHDKey.PublicKey().Serialize()

	// C.
	seed = []byte("this.is.c.seed.")
	cHDKey := bip32.NewHDKey(seed)
	cPrv := cHDKey.PrivateKey()
	cPub := cHDKey.PublicKey().Serialize()

	// Redeem script.
	redeemScript := xcore.NewPayToMultiSigScript(2, aPub, bPub, cPub)
	redeem, _ := redeemScript.GetLockingScriptBytes()
	fmt.Printf("redeem.hex:%x\n", redeem)
	multi := xcore.NewPayToWitnessV0ScriptHashAddress(xcrypto.Sha256(redeem))
	fmt.Printf("multi.addr:%v\n", multi.ToString(network.TestNet))

	// Funding to P2WSH.
	{
		bohuCoin := xcore.NewCoinBuilder().AddOutput(
			"092ddeb0fa8205a06494f2cf83afda0377479c86065e60dea5ae347468b27361",
			1,
			6758017,
			"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
		).ToCoins()[0]

		tx, err := xcore.NewTransactionBuilder().
			AddCoin(bohuCoin).
			AddKeys(bohuPrv).
			To(bohu, 4000).
			Then().
			SetChange(bohu).
			Then().
			Sign().
			BuildTransaction()
		assertNil(err)

		// Verify.
		err = tx.Verify()
		assertNil(err)

		fmt.Printf("p2wsh.fund:%v\n", tx.ToString())
		fmt.Printf("p2wsh.fund.txid:%v\n", tx.ID())
		fmt.Printf("p2wsh.fund.tx:%x\n", tx.Serialize())
	}

	// Spending From P2WSH.
	{
		multiCoin := xcore.NewCoinBuilder().AddOutput(
			"b2e955c95a6ee5752df1477a5936443ead0297ec697475ce6f356cdc6e2301a9",
			0,
			4000,
			"a914210a461ced66d7540ad2f4649b49dbed7c9fcc2887",
		).ToCoins()[0]

		tx, err := xcore.NewTransactionBuilder().
			AddCoin(multiCoin).
			AddKeys(aPrv, cPrv).
			SetRedeemScript(redeem).
			To(bohu, 1000).
			Then().
			SetChange(bohu).
			Then().
			Sign().
			BuildTransaction()
		assertNil(err)

		// Verify.
		err = tx.Verify()
		assertNil(err)

		fmt.Printf("p2wsh.spend:%v\n", tx.ToString())
		fmt.Printf("p2wsh.spend.txid:%v\n", tx.ID())
		fmt.Printf("p2wsh.spend.tx:%x\n", tx.Serialize())
	}
}
