# tokucore – A Simple, Powerful, Modular Library for Bitcoin and Blockchain-Based Apps.

[![Build Status](https://travis-ci.org/tokublock/tokucore.png)](https://travis-ci.org/tokublock/tokucore) [![Go Report Card](https://goreportcard.com/badge/github.com/tokublock/tokucore)](https://goreportcard.com/report/github.com/tokublock/tokucore) [![codecov.io](https://codecov.io/gh/tokublock/tokucore/graphs/badge.svg)](https://codecov.io/gh/tokublock/tokucore/branch/master)

## tokucore

*tokucore* is a simple Go (golang) library for creating and manipulating bitcoin data structures like creating keys and addresses (HD/bip32) or parsing, creating and signing transactions, micropayment.

## Focus

* simple and easy to use
* no external dependencies
* full test coverage

## Tests

```
$ export GOPATH=`pwd`
$ go get -u github.com/tokublock/tokucore/xcore
$ cd src/github.com/tokublock/tokucore/
$ make test
```
## Examples

### P2PKH

[p2pkh example](examples/p2pkh.go)

```
go run examples/p2pkh.go
```

### Multisig

[multisig example](examples/multisig.go)

```
go run examples/multisig.go
```

### HDWallet

[HDWallet example](examples/hdwallet.go)

```
go run examples/hdwallet.go
```

### Micropayment

[micropayment example](examples/micropayment.go)

```
go run examples/micropayment.go
```

## License

tokucore is released under the BSD.
