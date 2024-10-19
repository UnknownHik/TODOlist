package models

// Task описывает структуру входного запроса
type Task struct {
	ID      string `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment,omitempty" db:"comment"`
	Repeat  string `json:"repeat,omitempty" db:"repeat"`
}

// TaskResponse описывает структуру ответа
type TaskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// SignRequest содержит структуру для пароля из JSON-запроса
type Credentials struct {
	Password string `json:"password"`
}

// TaskFilter содержит структуру для поиска задачи по фильтру
type TaskFilter struct {
	Search     string
	SearchData bool
}

// JWTTokenResponse содержит структуру для ответа с токеном
type JWTTokenResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}
