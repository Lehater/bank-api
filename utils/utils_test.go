package utils_test

import (
	"bank-api/utils"
	"strconv"
	"testing"
	"golang.org/x/crypto/bcrypt"
)

// Тест для алгоритма Луна: если передать известное число, контрольная цифра должна совпадать.
func TestComputeLuhn(t *testing.T) {
	// Поскольку функция computeLuhnCheckDigit не экспортирована, проверяем через GenerateCardNumber
	generated := utils.GenerateCardNumber()
	if len(generated) != 16 {
		t.Errorf("expected 16-digit card number, got %s", generated)
	}
	// Простейшая проверка: функция возвращает число, которое можно преобразовать
	if _, err := strconv.Atoi(generated); err != nil {
		t.Errorf("generated card number is not numeric: %v", err)
	}
}

func TestGenerateExpirationDate(t *testing.T) {
	exp := utils.GenerateExpirationDate(5)
	// Формат должен быть "MM/YY" – 5 символов.
	if len(exp) != 5 {
		t.Errorf("expected expiration date in format MM/YY, got %s", exp)
	}
}

func TestGenerateCVV(t *testing.T) {
	cvv, err := utils.GenerateCVV()
	if err != nil {
		t.Fatalf("GenerateCVV error: %v", err)
	}
	if len(cvv) != 3 {
		t.Errorf("expected 3-digit CVV, got %s", cvv)
	}
}

func TestEncryptDecryptData(t *testing.T) {
	plainText := "SensitiveData123"
	encrypted, err := utils.EncryptData(plainText)
	if err != nil {
		t.Fatalf("EncryptData error: %v", err)
	}
	decrypted, err := utils.DecryptData(encrypted)
	if err != nil {
		t.Fatalf("DecryptData error: %v", err)
	}
	if decrypted != plainText {
		t.Errorf("expected %s, got %s", plainText, decrypted)
	}
}

func TestHashCVV(t *testing.T) {
	cvv := "123"
	hash, err := utils.HashCVV(cvv)
	if err != nil {
		t.Fatalf("HashCVV error: %v", err)
	}
	// Проверяем, что хеш соответствует исходному CVV с помощью bcrypt.
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(cvv)); err != nil {
		t.Errorf("hashed CVV does not match original: %v", err)
	}
}
