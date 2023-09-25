FROM golang:1.21-alpine AS build

RUN apk --no-cache add build-base

WORKDIR /app

COPY . /app

COPY go.mod go.sum ./

RUN go mod download

ENV CGO_ENABLED=1

RUN go build -tags musl -o ./cli ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/cli /bin/cli
COPY config.json /app/config.json

COPY --from=build /app/dao/postgres/migrations /app/dao/postgres/migrations

RUN ln -s /bin/cli /app/cli