package utils_test

import (
	"testing"

	"bank-api/utils"
)

func TestPGPEncryptDecrypt(t *testing.T) {
	plain := "4111111111111111"
	
	// Шифруем PGP + HMAC
	ciphertext, mac, err := utils.EncryptPGP(plain)
	if err != nil {
		t.Fatalf("EncryptPGP error: %v", err)
	}

	// Дешифруем и проверяем HMAC
	out, err := utils.DecryptPGP(ciphertext, mac)
	if err != nil {
		t.Fatalf("DecryptPGP error: %v", err)
	}

	if out != plain {
		t.Fatalf("Decrypted text mismatch: want %q, got %q", plain, out)
	}
}
