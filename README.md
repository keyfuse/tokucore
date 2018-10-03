# tokucore – A Simple, Powerful Library for Bitcoin Apps.

[![Build Status](https://travis-ci.org/tokublock/tokucore.png)](https://travis-ci.org/tokublock/tokucore) [![Go Report Card](https://goreportcard.com/badge/github.com/tokublock/tokucore)](https://goreportcard.com/report/github.com/tokublock/tokucore) [![codecov.io](https://codecov.io/gh/tokublock/tokucore/graphs/badge.svg)](https://codecov.io/gh/tokublock/tokucore/branch/master)

## tokucore

*tokucore* is a simple Go (golang) library for creating and manipulating bitcoin data structures like creating keys and addresses (HD/BIP32) or parsing, creating and signing transactions, micropayment.

## Focus

* Simple and easy to use
* No external dependencies
* Full test coverage

## Tests

```
$ export GOPATH=`pwd`
$ go get -u github.com/tokublock/tokucore/xcore
$ cd src/github.com/tokublock/tokucore/
$ make test
```

## Examples

- [P2PKH](examples/p2pkh.go)
- [2-to-3 MultiSig](examples/multisig.go)
- [HDWallet](examples/hdwallet.go)
- [MicroPayment](examples/micropayment.go)

## Applications

- [JustDoBlockchain](https://justdoblockchain.com) - A website Learning Blockchain Demo by Demo.

## License

tokucore is released under the BSD License.
