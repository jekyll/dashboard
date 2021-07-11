all: build test

build:
	go install ./...

test:
	go test ./...

server: build
	dashboard -http=localhost:8000
