.PHONY: build
build:
	go mod tidy
	go mod download
	go mod vendor
	go build -v ./cmd/zeronews

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build