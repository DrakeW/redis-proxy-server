.PHONY: build, run, stop, test

build:
	docker-compose build

run:
	docker-compose up -d

stop:
	docker-compose down

test: run
	go test ./...
