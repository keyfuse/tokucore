// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"testing"
)

func TestCoinBuilder(t *testing.T) {
	tests := []struct {
		name   string
		txid   string
		index  uint32
		script string
		value  uint64
	}{
		{
			name:   "#1",
			txid:   "0a4381c05136c0cb44886a5df7c26f1930bcc2c12e00ec60e027c4378d7d8c2e",
			index:  1,
			script: "a914203736c3c06053896d7041ce8f5bae3df76cc49187",
			value:  0.5 * 1e8,
		},
		{
			name:   "#2",
			txid:   "2c4df245d00b491bdf24965adbbccdaa7f62ccac933d3e9377f336c60c4ea096",
			index:  0,
			script: "a914f3ba8a120d960ae07d1dbe6f0c37fb4c926d76d587",
			value:  2.0 * 1e8,
		},
	}
	coinBuilder := NewCoinBuilder()
	for _, test := range tests {
		coinBuilder.AddOutput(test.txid, test.index, test.value, test.script)
	}
	coins := coinBuilder.ToCoins()
	t.Logf("coins:%v", coins)
}
