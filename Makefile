.SILENT:
.ONESHELL:
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.PHONY: run exec build_public build clean deps heroku_init heroku dev docker_build

NAME=$(shell basename $(CURDIR))

run: build_public build exec clean

exec:
	./bin/${NAME}

build_public:
	go-bindata -fs -pkg static -o pkg/static/static.go -prefix "static/" static

build:
	CGO_ENABLED=0 go build -trimpath -o bin/${NAME} -ldflags '-s -w -extldflags "-static"'

clean:
	rm -rf bin

deps:
	go mod init || true
	go mod tidy
	go mod verify

heroku_init:
	heroku login
	heroku git:clone -a draw-cool

heroku:
	git push heroku master

dev:
	go get -u -v github.com/go-bindata/go-bindata/...

docker_build:
	docker build -t chneau/draw:latest .
