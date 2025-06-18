# syntax=docker/dockerfile:1.7-labs
# Экспериментальцая версия для поддержки COPY --exclude
FROM golang:1.23.6-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY --exclude=go.mod --exclude=go.sum  . .

RUN CGO_ENABLED=0 GOOS=linux go build -o producer ./cmd/producer

FROM alpine:latest

COPY --from=builder /app/producer ./

ENTRYPOINT ["/producer"] 