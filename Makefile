
module: module.tar.gz

bin/viamoscmdmodule: go.mod *.go cmd/module/*.go
	go build -o bin/viamoscmdmodule cmd/module/cmd.go

lint:
	gofmt -s -w .

updaterdk:
	go get go.viam.com/rdk@latest
	go mod tidy

test:
	go test ./...


module.tar.gz: bin/viamoscmdmodule
	tar czf $@ $^

all: test bin/viamoscmd module 


