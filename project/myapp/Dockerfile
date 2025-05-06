# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Установка зависимостей
RUN apk add --no-cache git

# Копируем файлы модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем и собираем
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/group-service ./cmd/server

# Runtime stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/group-service .
COPY --from=builder /app/migrations ./migrations

# Настройка времени и пользователя
RUN apk add --no-cache tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
    echo "Europe/Moscow" > /etc/timezone && \
    adduser -D -g '' appuser && \
    chown -R appuser:appuser /app

USER appuser
EXPOSE 8080
CMD ["./group-service"]