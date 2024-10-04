package services

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) == 0 {
			next(w, r)
			return
		}

		// Получаем куку
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		// Проверяем валидности токена
		jwtToken, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TODO_JWT_SECRET")), nil
		})
		if err != nil || !jwtToken.Valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Генерируем хэш текущего пароля
		hash := sha256.New()
		hash.Write([]byte(pass))
		currentPassHash := hex.EncodeToString(hash.Sum(nil))

		// Получаем хэш из пароля токена
		passwordHash, ok := claims["passwordHash"].(string)
		if !ok {
			http.Error(w, "Failed to typecast to string", http.StatusUnauthorized)
			return
		}

		//Сравниваем текущий пароль с хэшем из токена
		if passwordHash != currentPassHash {
			http.Error(w, "Password has changed, please re-authenticate", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
