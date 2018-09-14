# build stage
ARG BASE=/go/src/github.com/chneau/draw
FROM golang:alpine AS build-env
ARG BASE
ADD . $BASE
RUN cd $BASE && CGO_ENABLED=0 go build -o /draw -ldflags '-s -w -extldflags "-static"'

FROM alpine AS prod-ready
COPY --from=build-env /draw /draw
ENTRYPOINT [ "/draw" ]
