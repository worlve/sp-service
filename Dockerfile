FROM golang:alpine

RUN mkdir -p /go/src/github.com/worlve/sp-service

WORKDIR /go/src/github.com/worlve/sp-service

CMD ["go", "run", "cmd/server/main.go"]