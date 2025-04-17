package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"bank-api/models"
)

// UserRepository определяет методы для работы с сущностью User.
type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id int) (*models.User, error)
}

// userRepository – конкретная реализация UserRepository.
type userRepository struct {
	db *sql.DB
}

// NewUserRepository возвращает новый экземпляр userRepository.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create добавляет нового пользователя в базу данных.
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, username, password_hash, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id
	`
	err := r.db.QueryRow(query, user.Email, user.Username, user.PasswordHash, user.CreatedAt).
		Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil
}

// GetByEmail находит пользователя по email.
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, username, password_hash, created_at
		FROM users WHERE email = $1
	`
	row := r.db.QueryRow(query, email)
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	return &user, nil
}

// GetByID находит пользователя по ID.
func (r *userRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, username, password_hash, created_at
		FROM users WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	return &user, nil
}
