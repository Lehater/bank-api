package config_test

import (
	"bank-api/config"
	"os"
	"testing"
)

func TestConnectDB_WithoutEnv(t *testing.T) {
	// Очистим тестовые переменные, если они установлены.
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")

	_, err := config.ConnectDB()
	if err == nil {
		t.Error("expected error when DB_* env vars not set, got nil")
	}
}
