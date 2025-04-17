package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"strconv"
	"time"
)

// Для демонстрации используется статический ключ AES-128 (16 байт).
var encryptionKey = []byte("example key 1234")

// GenerateCardNumber генерирует 16-значный номер карты, валидный по алгоритму Луна.
func GenerateCardNumber() string {
	number := ""
	// Генерируем 15 случайных цифр.
	for i := 0; i < 15; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		number += n.String()
	}
	// Вычисляем контрольную цифру по алгоритму Луна.
	checkDigit := computeLuhnCheckDigit(number)
	return number + strconv.Itoa(checkDigit)
}

func computeLuhnCheckDigit(number string) int {
	sum := 0
	double := true
	// Обходим цифры справа налево.
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}
	return (10 - (sum % 10)) % 10
}

// GenerateExpirationDate возвращает дату окончания срока действия карты в формате "MM/YY",
// смещенную на offsetYears лет от текущей даты.
func GenerateExpirationDate(offsetYears int) string {
	t := time.Now().AddDate(offsetYears, 0, 0)
	return fmt.Sprintf("%02d/%02d", t.Month(), t.Year()%100)
}

// GenerateCVV генерирует случайное трехзначное число в виде строки.
func GenerateCVV() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(n.Int64()) + 100), nil
}

// EncryptData шифрует переданный текст с помощью AES-GCM.
func EncryptData(plainText string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(ciphertext), nil
}

// DecryptData расшифровывает зашифрованный текст, полученный через EncryptData.
func DecryptData(encryptedHex string) (string, error) {
	ciphertext, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// HashCVV хеширует переданный CVV с использованием bcrypt.
func HashCVV(cvv string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
