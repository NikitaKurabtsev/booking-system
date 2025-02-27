# Build stage
FROM golang:1.22.4-alpine AS builder
WORKDIR /app

# Копируем зависимости и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем только API
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

# Устанавливаем migrate на этапе сборки
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Final stage
FROM alpine:latest
WORKDIR /app

# Копируем бинарник и зависимости
COPY --from=builder /app/bin/api .
COPY --from=builder /app/schema ./schema
COPY --from=builder /app/internal/config ./config

# Устанавливаем зависимости для здоровья
RUN apk add --no-cache postgresql-client curl

# Копируем migrate
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Порт для API
EXPOSE 8080

# Команда запуска
CMD ["./api"]