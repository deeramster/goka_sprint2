# Используем официальный образ Go
FROM golang:1.23

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Загружаем зависимости и собираем проект
WORKDIR /app/cmd/message-processor

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# Запускаем приложение
CMD ["./app"]
