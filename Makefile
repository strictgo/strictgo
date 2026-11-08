.PHONY: build install test vet fmt clean

build:
	go build -o bin/strictgo ./cmd/strictgo

install:
	go install ./cmd/strictgo

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

clean:
	rm -rf bin
