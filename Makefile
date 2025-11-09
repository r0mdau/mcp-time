.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

build: vet
	@echo ">> building mcp-time binary"
	mkdir -p build
	go build -o build/mcp-time ./cmd/mcp-time
.PHONY:build

test:
	go test -cover ./...
.PHONY:test

verify: fmt test
.PHONY:verify

run: build
	./build/mcp-time
.PHONY:run
