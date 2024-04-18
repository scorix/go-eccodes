# base stage for go
FROM golang:1.22-alpine3.19 as base

RUN apk add musl-dev gcc
RUN apk add eccodes --repository=https://dl-cdn.alpinelinux.org/alpine/edge/testing

WORKDIR /src
COPY . /src/

ENV CGO_ENABLED=1
RUN go build ./...
