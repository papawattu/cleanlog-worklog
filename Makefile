.PHONY: all build clean client run run-client watch docker-build

all: build client

build:
	@go build -o bin/cleanlog-worklog ./cmd/main.go

client:
	@go build -o bin/client ./client/main.go

run: build
	@./bin/cleanlog-worklog

run-client: client
	@./bin/client

clean:
	@rm -rf bin

watch:
	@air -c .air.toml

docker-build:
	@docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/papawattu/cleanlog-worklog:latest .