all: build test

deps:
	go get github.com/tools/godep

build: deps
	godep go build

test: deps
	godep go test
