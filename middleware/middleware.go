package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// UserIDKey используется для передачи идентификатора пользователя в контексте.
const UserIDKey = "userID"

// AuthMiddleware проверяет наличие и валидность JWT-токена.
// В случае успешной проверки извлекает идентификатор пользователя из токена 
// и добавляет его в контекст запроса.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Извлекаем данные из токена.
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || claims["sub"] == nil {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			// Добавляем идентификатор пользователя в контекст.
			userID, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Invalid user id in token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// LoggingMiddleware ведет логирование всех входящих HTTP-запросов,
// фиксируя метод, URI, время выполнения запроса и IP клиента.
func LoggingMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			logger.WithFields(logrus.Fields{
				"method":      r.Method,
				"path":        r.RequestURI,
				"duration":    duration,
				"remote_addr": r.RemoteAddr,
			}).Info("Handled request")
		})
	}
}

// RecoveryMiddleware отлавливает возможные паники (panic) в цепочке обработки запроса,
// логирует ошибку и возвращает статус 500 (Internal Server Error).
func RecoveryMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Errorf("Recovered from panic: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
