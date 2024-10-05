FROM golang:1.23.1

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
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o /todolist ./cmd/app/main.go

CMD ["/todolist"]