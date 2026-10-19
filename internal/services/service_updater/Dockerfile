FROM golang:1.26.2-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/updater_service \
    ./cmd/updater/main.go

FROM alpine:3.19
WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bin/updater_service /app/updater_service

CMD ["/app/updater_service"]