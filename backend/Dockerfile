FROM golang:1.22-alpine

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем приложение
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]