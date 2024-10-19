.PHONY: all build clean client run

all: build client

build:
	@go build -o bin/main *.go 

client:
	@go build -o bin/client client/main.go

run: build
	@./bin/main

run-client: client
	@./bin/client

clean:
	@rm -rf bin

watch:
	@air -c .air.toml