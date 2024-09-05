FROM golang:alpine AS stage
WORKDIR /go/src/github.com/safiachraf/DNServer-with-Redis
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o main main.go

FROM alpine:3.18

RUN apk update && apk add --no-cache redis

RUN apk add --no-cache redis

COPY --from=stage /go/src/github.com/kyhsa93/gin-rest-cqrs-example/main /main

EXPOSE 6379

CMD redis-server --daemonize no & /main
