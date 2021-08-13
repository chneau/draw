.SILENT:
.ONESHELL:
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.PHONY: run exec build clean deps heroku_init heroku dev docker_build

NAME=$(shell basename $(CURDIR))

run: build exec clean

exec:
	./bin/${NAME}

build:
	CGO_ENABLED=0 go build -trimpath -o bin/${NAME} -ldflags '-s -w -extldflags "-static"'

clean:
	rm -rf bin

deps:
	go mod init || true
	go mod tidy
	go mod verify

docker_build:
	docker build -t chneau/draw:latest .
