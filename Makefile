SHELL := /bin/bash

.PHONY: all clean install update

all: serve

install:
	go get -v ./...

dev-install: install
	cd .. && go install github.com/codegangsta/gin@latest

update:
	go get -v -u ./...

build: install
	go build ./cmd/jphotos-server

serve: install
	gin --build ./cmd/jphotos-server --all --excludeDir data -i

clean:
	rm jphotos-server gin-bin
