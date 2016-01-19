all: build test install

deps:
	go get github.com/tools/godep
	godep save ./... \
	  golang.org/x/net/context \
	  goji.io \
	  golang.org/x/oauth2 \
	  github.com/google/go-github/github

build: deps
	godep go build ./...

test: deps
	godep go test ./...

install: deps
	godep go install ./...
