package rest

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"todo-rest/internal/models"
)

// TokenHandler обрабатывает запросы на аутентификацию
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	var token models.JWTTokenResponse

	// Получаем пароль из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		token.Error = "Invalid request body"
		response(w, http.StatusBadRequest, token)
		return
	}

	// Получаем хранимый пароль из переменной окружения
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		token.Error = "Password not set in environment"
		response(w, http.StatusInternalServerError, token)
		return
	}

	// Проверяем, что введённый пароль совпадает с хранимым
	if creds.Password != pass {
		token.Error = "Wrong password"
		response(w, http.StatusUnauthorized, token)
		return
	}

	// Создаём JWT и указываем алгоритм хеширования
	passwordHash := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"passwordHash": passwordHash,
	})

	// Получаем подписанный токен
	signedToken, err := jwtToken.SignedString([]byte(os.Getenv("TODO_JWT_SECRET")))
	if err != nil {
		token.Error = "Failed to sign JWT"
		response(w, http.StatusBadRequest, token)
		return
	}

	response(w, http.StatusOK, models.JWTTokenResponse{Token: signedToken})
}
