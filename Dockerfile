# Build stage
FROM golang:1.22.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download -x && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w -X main.Version=$(git describe --tags --always)" \
    -trimpath \
    -o /app/bin/api ./cmd/api

# Final stage
FROM gcr.io/distroless/static-debian12:nonroot

LABEL maintainer="kurabtsevnikita@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/NikitaKurabtsev/booking-system"

# Environment variables
ENV ELASTICSEARCH_URL=http://elasticsearch:9200 \
    LOGSTASH_HOST=logstash:5000 \
    LOG_LEVEL=info

# Copy artifacts
WORKDIR /app
USER nonroot:nonroot
COPY --from=builder --chown=nonroot:nonroot /app/bin/api /app/api
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Network configuration
EXPOSE 80/tcp

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/app/api", "healthcheck"]

ENTRYPOINT ["/app/api"]