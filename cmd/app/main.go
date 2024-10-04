package main

import (
	"github.com/joho/godotenv"
	"log"
	"todo-rest/internal/database"
	"todo-rest/internal/transport/server"
)

func main() {
	// Инициализируем базу данных
	db := database.InitDb()
	defer db.Close()

	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Получаем порт и запускаем сервер
	port := server.GetPort()
	server.StartServer(port)
}
