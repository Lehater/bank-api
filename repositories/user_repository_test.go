package repositories_test

import (
	"bank-api/models"
	"bank-api/repositories"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserRepository_CreateAndGet(t *testing.T) {
	// Создаем мок для базы данных.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	// Подготавливаем тестового пользователя.
	user := &models.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
	}

	// Ожидаем запрос INSERT для создания пользователя.
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO users (email, username, password_hash, created_at) VALUES ($1, $2, $3, $4) RETURNING id`)).
		WithArgs(user.Email, user.Username, user.PasswordHash, user.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Вызываем Create.
	if err := repo.Create(user); err != nil {
		t.Errorf("unexpected error on Create: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("expected user ID 1, got %d", user.ID)
	}

	// Ожидаем SELECT для GetByEmail.
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, email, username, password_hash, created_at FROM users WHERE email = $1`)).
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "password_hash", "created_at"}).
			AddRow(1, user.Email, user.Username, user.PasswordHash, user.CreatedAt))

	gotUser, err := repo.GetByEmail(user.Email)
	if err != nil {
		t.Errorf("unexpected error on GetByEmail: %v", err)
	}
	if gotUser.ID != 1 {
		t.Errorf("expected user ID 1, got %d", gotUser.ID)
	}

	// Ожидаем SELECT для GetByID.
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, email, username, password_hash, created_at FROM users WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "password_hash", "created_at"}).
			AddRow(1, user.Email, user.Username, user.PasswordHash, user.CreatedAt))

	gotUser2, err := repo.GetByID(1)
	if err != nil {
		t.Errorf("unexpected error on GetByID: %v", err)
	}
	if gotUser2.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, gotUser2.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
