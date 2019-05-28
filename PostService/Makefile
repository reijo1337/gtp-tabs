deps:
	go test ./...
	go mod vendor

test:
	go test -mod=vendor -v ./...

build:
	CGO_ENABLED=0 go build -mod=vendor -o .bin/PostService *.go

run:build
	.bin/PostService