SHELL := /bin/bash
PWD := $(shell pwd)

clean:
	rm -rf */*.exe
.PHONY: clean

build:
	go build -v ./...
.PHONY: build

run:
	swag init
	go run .
.PHONY: run

test:
	go test -v ./...
.PHONY: test

docker-image:
	docker build -f ./Dockerfile -t "series-api:latest" .
	# Execute this command from time to time to clean up intermediate stages generated
	# during client build (your hard drive will like this :) ). Don't left uncommented if you
	# want to avoid rebuilding client image every time the docker-compose-up command
	# is executed, even when client code has not changed
	# docker rmi `docker images --filter label=intermediateStageToBeDeleted=true -q`
.PHONY: docker-image

