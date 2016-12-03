#.PHONY: build doc fmt lint dev test vet godep install bench
.PHONY: build test install

PKG_NAME=$(shell basename `pwd`)

install:
  - go get gopkg.in/redis.v5

build: test \
	go build -v -o ./bin/$(PKG_NAME)

test:
	go test ./...

