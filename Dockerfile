FROM golang:1.23.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники приложения в рабочую директорию
COPY . .

# Компилируем приложение
RUN go build -o user-service ./cmd/main.go

FROM alpine

WORKDIR /root/

# Копируем скомпилированное приложение из образа builder
COPY --from=builder /app/user-service .

COPY .env .

# Открываем порт 10101
EXPOSE 10101

# Запускаем приложение
CMD ["./user-service"]
