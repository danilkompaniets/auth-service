# Этап сборки
FROM golang:1.24.1-alpine AS build

WORKDIR /app

# Устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

RUN go mod tidy

# Копируем исходники
COPY . .

# Сборка бинарника auth-service
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd

# Финальный образ
FROM alpine:3.18

# Добавим сертификаты, чтобы можно было ходить в сеть (например, к БД через SSL)
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=build /app/auth-service .

EXPOSE 8081

CMD ["./auth-service"]