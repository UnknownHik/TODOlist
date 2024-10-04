package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"todo-rest/internal/services"

	"todo-rest/internal/config"
	"todo-rest/internal/database"
	"todo-rest/internal/models"
)

// NextDateHandler обрабатывает запросы к /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса
	now, err := time.Parse(config.DateFormat, r.FormValue("now"))
	if err != nil {
		http.Error(w, "Invalid now format. Expected format: 20060102", http.StatusInternalServerError)
		return
	}
	taskDate := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := services.NextDate(now, taskDate, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		log.Printf("Error writing response: %v\n", err)
	}
}

// response отправляет JSON ответ клиенту
func response(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// CreateTaskHandler обрабатывает POST запрос для добавления задачи
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var res models.TaskResponse

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		res.Error = "JSON deserialization error"
		response(w, http.StatusBadRequest, res)
		return
	}

	// Проверяем наличие заголовка
	if task.Title == "" {
		res.Error = "Task title not specified"
		response(w, http.StatusInternalServerError, res)
		return
	}

	// Получаем сегодняшнюю дату
	now := time.Now()

	// Проверяем формат даты и устанавливаем текущую дату
	if task.Date == "" {
		task.Date = now.Format(config.DateFormat)
	} else {
		date, err := time.Parse(config.DateFormat, task.Date)
		if err != nil {
			res.Error = "Date is in the wrong format"
			response(w, http.StatusInternalServerError, res)
			return
		}
		if date.Before(now) {
			task.Date = now.Format(config.DateFormat)
		}
	}

	// Проверяем правило повторения
	if task.Repeat != "" {
		if _, err := services.NextDate(now, task.Date, task.Repeat); err != nil {
			res.Error = "Invalid format of repeat rule"
			response(w, http.StatusInternalServerError, res)
			return
		}
	}

	// Добавляем задачу в базу данных
	taskId, err := database.AddTask(task)
	if err != nil {
		res.Error = "Failed to create task"
		response(w, http.StatusBadRequest, res)
		return
	}
	res.ID = fmt.Sprintf("%d", taskId)
	response(w, http.StatusOK, res)
}

// GetTasksListHandler обрабатывает GET запрос для вывода задач
func GetTasksListHandler(w http.ResponseWriter, r *http.Request) {
	var searchDate bool

	search := r.FormValue("search")
	searchParsed, err := time.Parse("02.01.2006", search)
	if err == nil {
		searchDate = true
		search = searchParsed.Format("20060102")
	} else {
		search = "%" + search + "%"
	}

	tasks, err := database.GetTasks(search, searchDate)
	if err != nil {
		res := models.TaskResponse{Error: "error getting task list"}
		response(w, http.StatusBadRequest, res)
		return
	}

	res := map[string]interface{}{"tasks": tasks}
	response(w, http.StatusOK, res)

}

// GetTaskIdHandler обрабатывает GET запрос для вывода параметров задачи
func GetTaskIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	task, err := database.GetTask(id)
	if err != nil {
		res := models.TaskResponse{Error: "failed to encode response"}
		response(w, http.StatusBadRequest, res)
		return
	}

	response(w, http.StatusOK, task)
}

// UpdateTaskHandler обрабатывает PUT запрос для обновления параметров задачи
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var res models.TaskResponse

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		res.Error = "JSON deserialization error"
		response(w, http.StatusBadRequest, res)
		return
	}

	// Проверяем наличие ID
	if task.ID == "" {
		res.Error = "Missing task ID"
		response(w, http.StatusInternalServerError, res)
		return
	}

	// Проверяем валидность ID
	if _, err := strconv.Atoi(task.ID); err != nil {
		res.Error = "Invalid ID"
		response(w, http.StatusInternalServerError, res)
		return
	}

	// Проверяем наличие заголовка
	if task.Title == "" {
		res.Error = "Task title not specified"
		response(w, http.StatusInternalServerError, res)
		return
	}

	// Получаем сегодняшнюю дату
	now := time.Now()

	// Проверяем формат даты и устанавливаем текущую дату
	dateParse, err := time.Parse(config.DateFormat, task.Date)
	if err != nil {
		res.Error = "Date is in the wrong format"
		response(w, http.StatusInternalServerError, res)
		return
	}
	if dateParse.Before(now) {
		task.Date = now.Format(config.DateFormat)
	}

	// Проверяем правило повторения
	if task.Repeat != "" {
		if _, err := services.NextDate(now, task.Date, task.Repeat); err != nil {
			res.Error = "Invalid format of repeat rule"
			response(w, http.StatusInternalServerError, res)
			return
		}
	}

	_, err = database.GetTask(task.ID)
	if err != nil {
		res := models.TaskResponse{Error: "Task not found"}
		response(w, http.StatusBadRequest, res)
		return
	}

	_, err = database.UpdateTask(task)
	if err != nil {
		response(w, http.StatusBadRequest, models.TaskResponse{Error: "Failed to update task"})
		return
	}

	response(w, http.StatusOK, task)
}

// DeleteTaskHandler обрабатывает DELETE запрос для удаления задачи
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	// Проверяем наличие ID
	if id == "" {
		response(w, http.StatusBadRequest, models.TaskResponse{Error: "Missing task ID"})
		return
	}

	// Проверяем валидность ID
	if _, err := strconv.Atoi(id); err != nil {
		response(w, http.StatusBadRequest, models.TaskResponse{Error: "Invalid ID"})
		return
	}

	if err := database.DeleteTask(id); err != nil {
		response(w, http.StatusBadRequest, models.TaskResponse{Error: "Failed to delete task"})
		return
	}

	response(w, http.StatusOK, struct{}{})
}

// DoneTaskHandler обрабатывает PUT запрос для отметки выполненных задач
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	task, err := database.GetTask(id)
	if err != nil {
		response(w, http.StatusBadRequest, models.TaskResponse{Error: "Task not found"})
		return
	}

	if task.Repeat == "" {
		if err := database.DeleteTask(task.ID); err != nil {
			response(w, http.StatusBadRequest, models.TaskResponse{Error: "Failed to delete task"})
			return
		}
	} else {
		task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			response(w, http.StatusInternalServerError, models.TaskResponse{Error: "Invalid format of repeat rule"})
			return
		}

		_, err := database.UpdateTask(task)
		if err != nil {
			response(w, http.StatusBadRequest, models.TaskResponse{Error: "Failed to update task"})
			return
		}
	}

	response(w, http.StatusOK, struct{}{})
}
