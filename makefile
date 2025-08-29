.PHONY: build test image all

test:
	go test -test.v ./cmd/...

build:
	go build -o bin/autoheal ./cmd

image:
	docker build -t autoheal .

all: test build image
