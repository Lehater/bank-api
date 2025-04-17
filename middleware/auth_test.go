package middleware_test

import (
	"bank-api/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

// generateTestJWT генерирует тестовый JWT с указанным userID.
func generateTestJWT(jwtSecret, userID string) string {
	claims := jwt.RegisteredClaims{
		Subject: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(jwtSecret))
	return tokenStr
}

func TestAuthMiddleware_Authorized(t *testing.T) {
	jwtSecret := "testsecret"
	tokenStr := generateTestJWT(jwtSecret, "42")

	// Handler, который просто возвращает 200 OK.
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.AuthMiddleware(jwtSecret)(finalHandler)
	req := httptest.NewRequest("GET", "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", rr.Code)
	}
}

func TestAuthMiddleware_Unauthorized(t *testing.T) {
	jwtSecret := "testsecret"

	// Handler, который возвращает 200 OK, если дошел до него.
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.AuthMiddleware(jwtSecret)(finalHandler)
	// Запрос без заголовка Authorization.
	req := httptest.NewRequest("GET", "/secure", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 Unauthorized, got %d", rr.Code)
	}

	// Проверка запроса с некорректным заголовком.
	req = httptest.NewRequest("GET", "/secure", nil)
	req.Header.Set("Authorization", "InvalidHeader")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 Unauthorized for malformed header, got %d", rr.Code)
	}
}
