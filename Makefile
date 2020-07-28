.SILENT:
.ONESHELL:
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.PHONY: run deps build clean exec dockerbuild

NAME=$(shell basename $(CURDIR))

run: buildPublic build exec clean

exec:
	./bin/${NAME}

buildPublic:
	go-bindata -fs -pkg static -o pkg/static/static.go -prefix "static/" static

build:
	CGO_ENABLED=0 go build -trimpath -o bin/${NAME} -ldflags '-s -w -extldflags "-static"'

clean:
	rm -rf bin

deps:
	go mod init || true
	go mod tidy
	go mod verify

dev:
	go get -u -v github.com/go-bindata/go-bindata/...

dockerbuild:
	docker build -t chneau/draw:latest .
