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

// Demo for funding coins to MultiSig address and spending from MultiSig address to P2PKH.
func main() {
	net := network.TestNet
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := bip32.NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohuPub := bohuHDKey.PublicKey()
	bohu := xcore.NewPayToPubKeyHashAddress(bohuPub.Hash160())
	fmt.Printf("bohu.addr:%v", bohu.ToString(net))

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
	multi := xcore.NewPayToScriptHashAddress(xcrypto.Hash160(redeem))
	fmt.Printf("multi.addr:%v\n", multi.ToString(net))

	// Funding.
	{
		bohuCoin := xcore.NewCoinBuilder().AddOutput(
			"b470aab9f18259b71fc7cb930929bedb6f6a15f7447219e7216db9a42c782984",
			0,
			129995000,
			"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
		).ToCoins()[0]

		tx, err := xcore.NewTransactionBuilder().
			AddCoin(bohuCoin).
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

		fmt.Printf("multisig.fund:%v\n", tx.ToString())
		fmt.Printf("multisig.fund.txid:%s\n", tx.ID())
		fmt.Printf("multisig.fund.tx:%x\n", tx.Serialize())
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

		fmt.Printf("multisig.spend:%v\n", tx.ToString())
		fmt.Printf("multisig.spend.txid:%s\n", tx.ID())
		fmt.Printf("multisig.spend.tx:%x\n", tx.Serialize())
	}
}
