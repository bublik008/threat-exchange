# --- Этап сборки (builder) ---
FROM golang:1.27-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Сначала копируем только файлы зависимостей, чтобы использовать кэш Docker
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник. 
# CGO_ENABLED=0 делает бинарник статическим (не зависит от системных библиотек)
RUN CGO_ENABLED=0 GOOS=linux go build -o /threat-exchange ./cmd/server

# --- Этап запуска (final) ---
# Используем пустой alpine-образ, чтобы контейнер был легким и безопасным
FROM alpine:3.19

# Устанавливаем сертификаты (на случай, если приложение будет ходить наружу по HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем скомпилированный бинарник из этапа сборки
COPY --from=builder /threat-exchange .

# Открываем порт (информативно, для документации)
EXPOSE 8080

# Команда по умолчанию
CMD ["./threat-exchange"]