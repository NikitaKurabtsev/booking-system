# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

ENV GOPROXY=https://goproxy.cn,direct

RUN go mod download
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Final stage
FROM alpine:3.19

RUN apk add --no-cache postgresql-client curl

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/internal/config ./config
COPY --from=builder /app/schema ./schema

EXPOSE 8080

CMD ["./api"]