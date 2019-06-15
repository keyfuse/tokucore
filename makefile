# tokucore
export PATH := $(GOPATH)/bin:$(PATH)

clean:
	@echo "--> Cleaning..."
	@go clean

test:
	go get -v github.com/stretchr/testify/assert
	@echo "--> Testing..."
	@$(MAKE) testxbase
	@$(MAKE) testxcrypto
	@$(MAKE) testxerror
	@$(MAKE) testxprotocol
	@$(MAKE) testnetwork
	@$(MAKE) testxvm
	@$(MAKE) testxcore

testxbase:
	go test -v -race ./xbase

testxcrypto:
	go test -v -race ./xcrypto/...

testxerror:
	go test -v -race ./xerror

testxrpc:
	go test -v -race ./xrpc

testxprotocol:
	go test -v -race ./xprotocol

testnetwork:
	go test -v -race ./network

testxvm:
	go test -v -race ./xvm

testxcore:
	go test -v -race ./xcore/bip32
	go test -v -race ./xcore/bip39
	go test -v -race ./xcore


pkgs =	./xbase\
		./xcrypto/...\
		./xerror\
		./xrpc\
		./xprotocol\
		./network\
		./xvm\
		./xcore/bip32\
		./xcore/bip39\
		./xcore

fmt:
	go vet $(pkgs)
	gofmt -s -w ./

coverage:
	go get -v github.com/pierrre/gotestcover
	gotestcover -coverprofile=coverage.out -v $(pkgs)
	go tool cover -html=coverage.out

check:
	go get -v gopkg.in/alecthomas/gometalinter.v2
	$(GOPATH)/bin/gometalinter.v2 -j 4 --disable-all \
	--enable=gofmt \
	--enable=golint \
	--enable=vet \
	--enable=gosimple \
	--enable=unconvert \
	--deadline=10m $(pkgs) 2>&1 | grep -v 'ALL_CAPS\|OP_' 2>&1 | tee /dev/stderr

.PHONY: clean test coverage check
