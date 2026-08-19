
module: module.tar.gz

bin/viam-system-module: go.mod *.go cmd/module/*.go
	go build -o bin/viam-system-module cmd/module/cmd.go

lint:
	gofmt -s -w .

updaterdk:
	go get go.viam.com/rdk@latest
	go mod tidy

test:
	go test ./...


module.tar.gz: bin/viam-system-module
	tar czf $@ $^

all: test bin/viam-system-module module

