deps:
	go get -u ./...
	go mod tidy
	go mod vendor

test:
	go test -mod=vendor -v ./...

build:
	CGO_ENABLED=0 go build -mod=vendor -o .bin/ProfileService *.go

run:build
	.bin/ProfileService