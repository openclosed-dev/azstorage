.DEFAULT_GOAL := build

.PHONY: build clean

build:
	go build -trimpath -o bin/ ./cmd/azstorage

clean:
	@rm -rf bin
