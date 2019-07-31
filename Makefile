.PHONY: build, run, stop, test

build:
	go build -o ./server main.go

run:
	docker-compose up -d

stop:
	docker-compose down

test: run
	go test ./...
