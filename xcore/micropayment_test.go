// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tokublock/tokucore/network"
)

func TestMicroPayment(t *testing.T) {
	var fee uint64 = 0.00001 * Unit
	var amount uint64 = 0.01 * Unit
	var cupOfCoffee int = 0.0001 * Unit

	var satoshiChannel *MicroPayer
	var starbucksChannel *MicroPayee

	locktime := uint32(time.Now().Add(time.Hour).Unix())

	satoshiCoin := NewCoinBuilder().AddOutput(
		"2b307391ccab9de7b91953cd488f35ddcd5bbc0c217cc0433009451a634d1db4",
		0,
		14854725,
		"76a9148b7f2212ecc4384abcf1df3fc5783e9c2a24d5a588ac",
	).ToCoins()[0]

	// Step#1. Create Channels.
	{
		// satoshi.
		seed := []byte("this.is.satoshi.seed.")
		satoshiHDKey := NewHDKey(seed)
		satoshiPrv := satoshiHDKey.PrivateKey()
		satoshiPub := satoshiHDKey.PublicKey()

		// starbucks.
		seed = []byte("this.is.starbucks.seed.")
		starbucksHDKey := NewHDKey(seed)
		starbucksPrv := starbucksHDKey.PrivateKey()
		starbucksPub := starbucksHDKey.PublicKey()

		satoshiChannel = NewMicroPayer(satoshiPrv, starbucksPub, amount, fee, locktime)
		t.Logf("satoshi.addr:%v", satoshiChannel.Address().ToString(network.TestNet))

		starbucksChannel = NewMicroPayee(starbucksPrv, satoshiPub, amount, fee, locktime)
		t.Logf("starbucks.addr:%v", starbucksChannel.Address().ToString(network.TestNet))
	}

	// Step#2. Channels handshake.
	{
		// #1.1 satoshi bond transaction.
		bondTx, err := satoshiChannel.CreateBond(satoshiCoin)
		assert.Nil(t, err)

		// #1.2 satoshi refund transaction.
		refundUnsignedTx, err := satoshiChannel.CreateRefund(bondTx)
		assert.Nil(t, err)

		// #2. starbucks gets satoshi's refund tx by tx.SerializeForPartially/DeserializeForPartially and signs it.
		starbucksRefundSign, err := starbucksChannel.SignRefund(refundUnsignedTx)
		assert.Nil(t, err)

		// #3. satoshi gets and verify starbucks's sign and sings the refund.
		refundTx, err := satoshiChannel.SignRefund(refundUnsignedTx, starbucksRefundSign)
		assert.Nil(t, err)

		// #4. starbucks check the refund and bond transaction.
		err = starbucksChannel.CheckBond(refundTx, bondTx)
		assert.Nil(t, err)

		t.Logf("bond.tx:%s", bondTx.ToString())
		t.Logf("bond.hex:%x", bondTx.Serialize())
		t.Logf("refund.hex:%x", refundTx.Serialize())
	}

	// Payments.
	// satoshi buys 4 cups coffee, then the channel closed by starbucks.
	for i := 1; i < 5; i++ {
		amount := i * cupOfCoffee
		paymentSign, err := satoshiChannel.SignPayment(uint64(amount))
		assert.Nil(t, err)

		tx, err := starbucksChannel.SignPayment(uint64(amount), paymentSign)
		assert.Nil(t, err)
		t.Logf("payment#%d.tx:%v", i, tx.ToString())
		t.Logf("payment#%d.hex:%x", i, tx.Serialize())
	}
	// starbucks brodcast the last payment and close the channel.
}
