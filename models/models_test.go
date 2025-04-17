package models_test

import (
	"bank-api/models"
	"encoding/json"
	"testing"
	"time"
)

func TestUserJSON(t *testing.T) {
	user := models.User{
		ID:        1,
		Email:     "test@example.com",
		Username:  "testuser",
		CreatedAt: time.Now().UTC(),
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal User: %v", err)
	}

	var user2 models.User
	err = json.Unmarshal(data, &user2)
	if err != nil {
		t.Fatalf("failed to unmarshal User: %v", err)
	}
	if user.Email != user2.Email || user.Username != user2.Username {
		t.Errorf("expected %v, got %v", user, user2)
	}
}
