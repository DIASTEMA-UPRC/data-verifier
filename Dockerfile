FROM golang:1.19-alpine AS builder

RUN apk update
RUN apk upgrade --available

WORKDIR /app

ADD . .

RUN go mod download
RUN go build -o main .

FROM alpine:latest

RUN apk update
RUN apk upgrade --available

WORKDIR /app

COPY --from=builder /app/main .

ENTRYPOINT ./main
