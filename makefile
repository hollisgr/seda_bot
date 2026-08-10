.PHONY: get

deps:
	go get -u ./...
	go mod tidy