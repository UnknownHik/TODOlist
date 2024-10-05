# TODOlist

TODOlist — это веб-приложение для управления задачами, разработанное на Go с использованием SQLite в качестве базы данных. Это приложение позволяет пользователям создавать, редактировать и удалять задачи.

## Установка

### Предварительные требования

- Go (версия 1.23 или выше)
- SQLite
- Docker

### Клонирование репозитория

```bash
git clone https://github.com/UnknownHik/TODOlist.git
cd TODOlist
```

## Запуск приложения

### Установите зависимости:
```bash
go mod tidy
```

### Настройте базу данных:  
Приложение использует SQLite. 
Убедитесь, что файл базы данных scheduler.db находится в корневом каталоге проекта. 
Если файл отсутствует, он будет создан автоматически при первом запуске приложения.  

### Запустите приложение:  
```bash
go run ./cmd/app/main.go
```
Приложение будет доступно по адресу http://localhost:7540.

## Запуск с использованием Docker
Если вы хотите запустить приложение с использованием Docker, выполните следующие шаги:

### Постройте Docker-образ:
```bash
docker build -t todolist .
```

### Запустите контейнер:
```bash
docker run -p 7540:7540 todolist
```

## Настройка

### config
Приложение может быть настроено через файл конфигурации. 
Убедитесь, что в папке internal/config есть файл конфигурации с правильными параметрами.

### .env
Приложение использует переменные окружения:
```bash
TODO_PORT - Порт, на котором будет запущено приложение (по умолчанию: 7540)
TODO_DBFILE - Путь к файлу базы данных	(по умолчанию: ./scheduler.db)
TODO_PASSWORD - Пароль для доступа (по умолчанию: aaa)
TODO_JWT_SECRET - Секретный ключ для подписи токена JWT (по умолчанию: secret)
```

## Тестирование
Для запуска тестов выполните следующую команду в корневом каталоге проекта:
```bash
go test ./tests
```
Это запустит все тесты, находящиеся в проекте, и выведет результаты на консоль.

## Аутентификация
Для управления доступом используется JWT (JSON Web Token). 
Перед выполнением операций, требующих аутентификации, убедитесь, что вы получили токен доступа.