.DEFAULT_GOAL := build

.PHONY: build build-win clean

build:
	go build -trimpath -o bin/ ./cmd/azstorage

build-win:
	GOOS=windows GOARCH=amd64 go build -trimpath -o bin/ ./cmd/azstorage

clean:
	@rm -rf bin
