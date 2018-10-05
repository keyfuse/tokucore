// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package main

import (
	"fmt"

	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xcore"
	"github.com/tokublock/tokucore/xvm"
)

func main() {
	// x+6=7
	hash160, err := xvm.NewScriptBuilder().
		AddOp(xvm.OP_6).
		AddOp(xvm.OP_ADD).
		AddOp(xvm.OP_7).
		AddOp(xvm.OP_EQUAL).Hash160()
	if err != nil {
		panic(err)
	}
	addr := xcore.NewPayToScriptHashAddress(hash160)

	fmt.Printf("p2sh.address(mainet):\t%s\n", addr.ToString(network.MainNet))
}
