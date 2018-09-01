# tokucore
export PATH := $(GOPATH)/bin:$(PATH)

clean:
	@echo "--> Cleaning..."
	@go clean

test:
	go get github.com/stretchr/testify/assert
	@echo "--> Testing..."
	@$(MAKE) testxbase
	@$(MAKE) testxcrypto
	@$(MAKE) testxerror
	@$(MAKE) testxrpc
	@$(MAKE) testxvm
	@$(MAKE) testxcore

testxbase:
	go test -v -race ./xbase/base58

testxcrypto:
	go test -v -race ./xcrypto
	go test -v -race ./xcrypto/ripemd160

testxerror:
	go test -v -race ./xerror

testxrpc:
	go test -v -race ./xrpc

testxvm:
	go test -v -race ./xvm

testxcore:
	go test -v -race ./xcore

bench:
	go test -run - -bench . ./...

fmt:
	go vet ./...
	gofmt -s -w ./

pkgs =	./xbase/base58\
		./xcrypto\
		./xerror\
		./xrpc\
		./xvm\
		./xcore

coverage:
	go get github.com/pierrre/gotestcover
	gotestcover -coverprofile=coverage.out -v $(pkgs)
	go tool cover -html=coverage.out

check:
	go get gopkg.in/alecthomas/gometalinter.v2
	go get honnef.co/go/tools/cmd/megacheck
	bin/gometalinter.v2 -j 4 --disable-all \
	--enable=gofmt \
	--enable=golint \
	--enable=vet \
	--enable=gosimple \
	--enable=unconvert \
	--deadline=10m $(pkgs) 2>&1 | grep -v 'ALL_CAPS\|OP_' 2>&1 | tee /dev/stderr
	bin/megacheck $(pkgs) 2>&1 | tee /dev/stderr

.PHONY: clean test coverage check
