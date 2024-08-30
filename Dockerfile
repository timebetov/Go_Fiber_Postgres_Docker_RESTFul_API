FROM golang:1.20-alpine

RUN apk add --no-cache make

WORKDIR /app

COPY ./src/go.mod ./src/go.sum ./

RUN go install github.com/cespare/reflex@latest
RUN go mod download

COPY ./src .
