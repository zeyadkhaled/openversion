FROM golang:1.14.2-alpine3.11 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --no-cache git

WORKDIR /backend

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build ./cmd/backend

FROM alpine:3.11

RUN apk add --no-cache

WORKDIR /app

COPY --from=builder /backend/backend /app

ENTRYPOINT ["/app/backend"]