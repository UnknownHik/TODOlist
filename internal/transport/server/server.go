package server

import (
	"log"
	"net/http"
	"os"

	"todo-rest/internal/config"
	"todo-rest/internal/services"
	"todo-rest/internal/transport/rest"

	"github.com/go-chi/chi/v5"
)

// GetPort получает порт из переменной окружения TODO_PORT или использует 7540 по умолчанию
func GetPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return port
}

// StartServer запускает веб-сервер на указанном порту
func StartServer(port string) {
	r := chi.NewRouter()

	cfg := config.LoadJWTConfig()

	// Настраиваем файловый сервер для каталога ./web
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))

	// Регистрация обработчика для API
	r.Post("/api/signin", rest.TokenHandler)
	r.Route("/api", func(r chi.Router) {
		r.Get("/nextdate", rest.NextDateHandler)
		r.Post("/task", services.Auth(cfg, rest.CreateTaskHandler))
		r.Get("/task", services.Auth(cfg, rest.GetTaskIdHandler))
		r.Put("/task", services.Auth(cfg, rest.UpdateTaskHandler))
		r.Delete("/task", services.Auth(cfg, rest.DeleteTaskHandler))
		r.Post("/task/done", services.Auth(cfg, rest.DoneTaskHandler))
		r.Get("/tasks", services.Auth(cfg, rest.GetTasksListHandler))
	})

	log.Printf("Server is running on port: %s\n", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
