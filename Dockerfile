FROM golang:1.14.2-alpine3.11 AS builder

RUN apk add --no-cache git

WORKDIR /opentelemetry_demo

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPRIVATE=gitlab.innology.com.tr

COPY . .

ENTRYPOINT [ "go", "run" ,"./cmd"]