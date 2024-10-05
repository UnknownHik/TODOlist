package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	"todo-rest/internal/config"
	"todo-rest/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDb инициализирует базу данных и создаёт необходимые таблицы и индексы, если они не существуют
func InitDb() *sql.DB {
	// Получаем путь к базе данных из переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		appPath, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dbFile = filepath.Join(appPath, "scheduler.db")
	}
	log.Println("Path to the database: ", dbFile)

	// Проверка существования базы данных
	_, err := os.Stat(dbFile)
	if err != nil && os.IsNotExist(err) {
		// Файл базы данных не существует, создаем новый
		file, err := os.Create(dbFile)
		if err != nil {
			log.Fatalf("Error creating database: %v", err)
		}
		file.Close()
	}

	// Открываем или создаем базу данных
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Создаем таблицу и индекс
	createTable := `
    CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date VARCHAR(10) NOT NULL,
        title VARCHAR(128) NOT NULL,
        comment TEXT,
        repeat VARCHAR(128)
    );
    CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
    `
	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	return db
}

// AddTask добавляет задачу в базу данных
func AddTask(task models.Task) (int, error) {
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetTasks выводит список всех задач или по фильтру
func GetTasks(filter models.TaskFilter) (tasks []models.Task, err error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit"
	if filter.Search != "" && !filter.SearchData {
		query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit"
	} else if filter.Search != "" && filter.SearchData {
		query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :search LIMIT :limit"
	}
	rows, err := db.Query(query, sql.Named("search", filter.Search), sql.Named("limit", config.LimitSearch))
	if err != nil {
		return []models.Task{}, errors.New("error getting task list")
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []models.Task{}, errors.New("data reading error")
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return []models.Task{}, errors.New("data reading error")
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

func GetTask(id string) (models.Task, error) {
	var task models.Task

	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return models.Task{}, err
	}

	return task, nil
}

// UpdateTask изменяет параметры задачи
func UpdateTask(task models.Task) (models.Task, error) {
	res, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return models.Task{}, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Task{}, err
	}

	if rowsAffected == 0 {
		return models.Task{}, errors.New("failed to delete")
	}

	return task, nil

}

// DeleteTask удаляет задачу
func DeleteTask(id string) error {
	res, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("failed to delete")
	}

	return nil

}
