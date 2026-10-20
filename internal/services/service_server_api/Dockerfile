FROM golang:1.26.2-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/server_api \
    ./cmd/reader/main.go

FROM alpine:3.19
WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bin/server_api /app/server_api

CMD ["/app/server_api"]