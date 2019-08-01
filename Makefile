.PHONY: dep, build, run, stop, test

clean:
	rm -rf ./vendor

dep: clean
	command -v dep >/dev/null || curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep ensure

build: dep
	go build -o ./server main.go

run:
	docker-compose up -d

stop:
	docker-compose down

test: run
	go test ./...
