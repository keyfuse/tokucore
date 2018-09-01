// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

// Coin -- the coin output from the pre tx.
type Coin struct {
	txID   string
	n      uint32
	value  uint64
	script string
}

// CoinBuilder --
type CoinBuilder struct {
	coins []*Coin
}

// NewCoinBuilder -- creates new CoinBuilder
func NewCoinBuilder() *CoinBuilder {
	return &CoinBuilder{}
}

// AddOutput -- add output to coin builder.
func (b *CoinBuilder) AddOutput(txID string, n uint32, value uint64, script string) *CoinBuilder {
	coin := &Coin{
		txID:   txID,
		n:      n,
		value:  value,
		script: script,
	}
	b.coins = append(b.coins, coin)
	return b
}

// ToCoins -- returns the coin slice.
func (b *CoinBuilder) ToCoins() []*Coin {
	return b.coins
}
