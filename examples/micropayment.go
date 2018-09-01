// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package main

import (
	"fmt"
	"time"

	"github.com/tokublock/tokucore/xcore"
)

func assertNil(err error) {
	if err != nil {
		panic(err)
	}
}

// Bitcoin MicroPayment channel.
func main() {
	var fee uint64 = 0.00001 * xcore.Unit
	var amount uint64 = 1 * xcore.Unit
	var cupOfCoffee int = 0.0001 * xcore.Unit

	var satoshiChannel *xcore.MicroPayer
	var starbucksChannel *xcore.MicroPayee

	locktime := uint32(time.Now().Add(time.Hour).Unix())

	satoshiCoin := xcore.NewCoinBuilder().AddOutput(
		"1f6f5669fb742147c5fa5fbde99f8cb21c65e6ff09bf045bc7d934f357b8fe15",
		0,
		2*xcore.Unit,
		"76a9148b7f2212ecc4384abcf1df3fc5783e9c2a24d5a588ac",
	).ToCoins()[0]

	// Step#1. Create Channels.
	{
		// satoshi.
		seed := []byte("this.is.satoshi.seed.")
		satoshiHDKey := xcore.NewHDKey(seed)
		satoshiPrv := satoshiHDKey.PrivateKey()
		satoshiPub := satoshiHDKey.PublicKey()

		// starbucks.
		seed = []byte("this.is.starbucks.seed.")
		starbucksHDKey := xcore.NewHDKey(seed)
		starbucksPrv := starbucksHDKey.PrivateKey()
		starbucksPub := starbucksHDKey.PublicKey()

		satoshiChannel = xcore.NewMicroPayer(satoshiPrv, starbucksPub, amount, fee, locktime)
		starbucksChannel = xcore.NewMicroPayee(starbucksPrv, satoshiPub, amount, fee, locktime)
	}

	// Step#2. Channels handshake.
	{
		// #1.1 satoshi bond transaction.
		bondTx, err := satoshiChannel.CreateBond(satoshiCoin)
		assertNil(err)

		// #1.2 satoshi refund transaction.
		refundUnsignedTx, err := satoshiChannel.CreateRefund(bondTx)
		assertNil(err)

		// #2. starbucks gets satoshi's refund tx by tx.SerializeForPartially/DeserializeForPartially and signs it.
		starbucksRefundSign, err := starbucksChannel.SignRefund(refundUnsignedTx)
		assertNil(err)

		// #3. satoshi gets and verify starbucks's sign and sings the refund.
		refundTx, err := satoshiChannel.SignRefund(refundUnsignedTx, starbucksRefundSign)
		assertNil(err)

		// #4. starbucks checks the refund and bond transaction.
		err = starbucksChannel.CheckBond(refundTx, bondTx)
		assertNil(err)

		fmt.Printf("bond.hex:%x\n", bondTx.Serialize())
		fmt.Printf("refund.hex:%x\n", refundTx.Serialize())
	}

	// Step#3. Payments.
	// satoshi buys 4 cups coffee, then the channel closed by starbucks.
	for i := 1; i < 5; i++ {
		amount := i * cupOfCoffee
		paymentSign, err := satoshiChannel.SignPayment(uint64(amount))
		assertNil(err)

		tx, err := starbucksChannel.SignPayment(uint64(amount), paymentSign)
		assertNil(err)
		fmt.Printf("payment#%d.tx:%s\n", i, tx.ToString())
		fmt.Printf("payment#%d.hex:%x\n", i, tx.Serialize())
	}
	// starbucks brodcast the last payment and close the channel.
}
