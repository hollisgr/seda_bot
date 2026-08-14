BIN=app
BUILD=cmd/main.go

.PHONY: get build run

all: build run

build:
	go build -o $(BIN) $(BUILD)

run:
	./$(BIN)

docker-compose-up-silent: docker-compose-stop
	sudo docker compose -f docker-compose.yml --env-file=.env up -d --build

docker-compose-stop:
	sudo docker compose -f docker-compose.yml --env-file=.env stop

docker-compose-up: docker-compose-down
	sudo docker compose -f docker-compose.yml --env-file=.env up --build

docker-compose-down:
	sudo docker compose -f docker-compose.yml --env-file=.env down

deps:
	go get -u ./...
	go mod tidy