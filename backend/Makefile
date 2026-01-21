

all: help

build: deps
	go build -o guestbook-api .

test:
	go test -v .

verbose:
	go test -v -race .

coverage:
	go test -coverprofile=coverage.out .
	go tool cover -html=coverage.out -o coverage.html

smoke:
	go test -short .

bench:
	go test -bench=. -benchmem .

deps:
	go mod download
	go mod tidy

clean:
	rm -f guestbook-api
	rm -f coverage.out coverage.html

help:
	@echo "make [test|verbose|coverage|smoke|bench]"
	@echo "make deps  - install dependencies"
	@echo "make clean - remove files, clean up dir"
	@echo "make help  - get help"



