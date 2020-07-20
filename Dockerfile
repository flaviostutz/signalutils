FROM golang:1.14.0-alpine3.11

RUN apk add --no-cache build-base

WORKDIR /signalutils

ADD go.mod .
RUN go mod download

ADD / /signalutils
RUN go test -v

WORKDIR /signalutils/example
RUN go build -o /usr/bin/signalutils-example

CMD [ "signalutils-example" ]
