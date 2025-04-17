# 1. Этап сборки (builder)
FROM golang:1.23-alpine AS builder

# 2. Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# 3. Копируем файлы зависимостей
#COPY ./go.mod ./go.sum 

# 4. Загружаем зависимости
#RUN go mod download

# 5. Копируем весь проект
COPY . .

# 6. Собираем бинарный файл приложения
#RUN go build -o task-manager ./cmd/server
RUN go mod tidy && go build -o task-manager ./cmd/server
# ----------------------------
# 7. Этап запуска (финальный образ)
FROM alpine:latest

# 8. Устанавливаем рабочую директорию
WORKDIR /app

# 9. Копируем бинарник из этапа builder
COPY --from=builder /app/task-manager .

# 10. Копируем миграции (если есть) и скрипты
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/scripts ./scripts

# Устанавливаем клиент PostgreSQL
RUN apk add --no-cache postgresql-client

# Делаем скрипт исполняемым
RUN chmod +x ./scripts/migrate.sh

# 11. Открываем порт для доступа извне
EXPOSE 8080

# Запускаем миграции перед стартом приложения
CMD ["./scripts/migrate.sh"]

# 12. Команда запуска приложения
CMD ["./task-manager"]