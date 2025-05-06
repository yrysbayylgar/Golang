package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	httputil "myapp/pkg/http"
)

// userIDKey - ключ контекста для ID пользователя
type userIDKey struct{}

// Logger - промежуточное ПО, которое логирует запросы
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Начат %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Завершен %s %s за %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// JSONContentType - промежуточное ПО, которое устанавливает тип содержимого application/json
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// JWTAuth - промежуточное ПО, которое проверяет JWT токены
func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Извлекаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httputil.RespondWithError(w, http.StatusUnauthorized, "Требуется заголовок Authorization")
				return
			}

			// Токен должен быть в формате "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				httputil.RespondWithError(w, http.StatusUnauthorized, "Недопустимый формат авторизации")
				return
			}

			tokenString := parts[1]

			// Парсим и проверяем токен
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Проверяем алгоритм
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				httputil.RespondWithError(w, http.StatusUnauthorized, "Недействительный токен")
				return
			}

			// Извлекаем ID пользователя из утверждений (claims)
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				httputil.RespondWithError(w, http.StatusUnauthorized, "Недействительные утверждения токена")
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				httputil.RespondWithError(w, http.StatusUnauthorized, "Недействительный ID пользователя в токене")
				return
			}

			// Добавляем ID пользователя в контекст
			ctx := context.WithValue(r.Context(), userIDKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID извлекает ID пользователя из контекста
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey{}).(string)
	return userID, ok
}