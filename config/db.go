package config

import (
    "database/sql"
    "fmt"
    "os"

    _ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || password == "" || name == "" || port == "" {
		return nil, fmt.Errorf("missing database environment variables")
	}

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

    return sql.Open("postgres", dsn)
}

