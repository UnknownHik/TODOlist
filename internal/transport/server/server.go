package server

import (
	"log"
	"net/http"
	"os"
	"todo-rest/internal/services"

	"github.com/go-chi/chi/v5"
	"todo-rest/internal/transport/rest"
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

	// Настраиваем файловый сервер для каталога ./web
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))

	// Регистрация обработчика для API
	r.Post("/api/signin", rest.TokenHandler)
	r.Route("/api", func(r chi.Router) {
		r.Get("/nextdate", rest.NextDateHandler)
		r.Post("/task", services.Auth(rest.CreateTaskHandler))
		r.Get("/task", services.Auth(rest.GetTaskIdHandler))
		r.Put("/task", services.Auth(rest.UpdateTaskHandler))
		r.Delete("/task", services.Auth(rest.DeleteTaskHandler))
		r.Post("/task/done", services.Auth(rest.DoneTaskHandler))
		r.Get("/tasks", services.Auth(rest.GetTasksListHandler))
	})

	log.Printf("Server is running on port: %s\n", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
