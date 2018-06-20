FROM golang:latest
WORKDIR /go/src/github.com/chneau/draw
COPY . .

ENV CGO_ENABLED=0
RUN go get -v ./...
RUN go build -o /app -ldflags '-s -w -extldflags "-static"'

FROM chneau/upx:latest
COPY --from=0 /app /app
RUN upx /app

FROM scratch
COPY --from=1 /app /app
ENTRYPOINT ["/app"]
