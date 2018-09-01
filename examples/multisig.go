// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package main

import (
	"fmt"

	"github.com/tokublock/tokucore/xcore"
)

func assertNil(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := xcore.NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohu := bohuHDKey.GetAddress()
	fmt.Printf("bohu.addr:%v", bohu.ToString(xcore.TestNet))

	// A.
	seed = []byte("this.is.a.seed.")
	aHDKey := xcore.NewHDKey(seed)
	aPrv := aHDKey.PrivateKey()
	aPub := aHDKey.PublicKey().Serialize()

	// B.
	seed = []byte("this.is.b.seed.")
	bHDKey := xcore.NewHDKey(seed)
	bPub := bHDKey.PublicKey().Serialize()

	// C.
	seed = []byte("this.is.c.seed.")
	cHDKey := xcore.NewHDKey(seed)
	cPrv := cHDKey.PrivateKey()
	cPub := cHDKey.PublicKey().Serialize()

	// Redeem script.
	redeemScript := xcore.NewPayToMultiSigScript(2, aPub, bPub, cPub)
	multi := redeemScript.GetAddress()
	fmt.Printf("multi.addr:%v\n", multi.ToString(xcore.TestNet))
	redeem, _ := redeemScript.GetLockingScriptBytes()
	fmt.Printf("redeem.hex:%x\n", redeem)

	// Funding.
	{
		bohuCoin := xcore.NewCoinBuilder().AddOutput(
			"b470aab9f18259b71fc7cb930929bedb6f6a15f7447219e7216db9a42c782984",
			0,
			129995000,
			"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
		).ToCoins()[0]

		tx, err := xcore.NewTransactionBuilder().
			AddCoins(bohuCoin).
			AddKeys(bohuPrv).
			To(multi, 4000).
			Then().
			SetChange(bohu).
			Then().
			Sign().
			BuildTransaction()
		assertNil(err)

		// Verify.
		err = tx.Verify()
		assertNil(err)

		fmt.Printf("%v\n", tx.ToString())
		fmt.Printf("funding.txid:%s\n", tx.ID())
		signedTx := tx.Serialize()
		fmt.Printf("funding.signed.tx:%x\n", signedTx)
	}

	// Spending.
	{
		multiCoin := xcore.NewCoinBuilder().AddOutput(
			"5af1520f1d3e818fca695c2a903baa4a7eec4954f0b35aa01be1f2c1d2cfd802",
			1,
			4000,
			"a914210a461ced66d7540ad2f4649b49dbed7c9fcc2887",
		).ToCoins()[0]

		tx, err := xcore.NewTransactionBuilder().
			AddCoins(multiCoin).
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

		fmt.Printf("%v\n", tx.ToString())
		fmt.Printf("spending.txid:%s\n", tx.ID())
		signedTx := tx.Serialize()
		fmt.Printf("spending.signed.tx:%x\n", signedTx)
	}
}
