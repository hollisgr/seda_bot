BIN=app
BUILD=cmd/main.go

.PHONY: get build run

all: build run

build:
	go build -o $(BIN) $(BUILD)

run:
	./$(BIN)

deps:
	go get -u ./...
	go mod tidy