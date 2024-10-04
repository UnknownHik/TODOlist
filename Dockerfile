FROM golang:1.23

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости Go
RUN go mod download

# Устанавливаем переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=./scheduler.db
ENV TODO_PASSWORD=aaa
ENV TODO_JWT_SECRET=secret

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /todolist ./cmd/app/main.go

CMD ["/todolist"]