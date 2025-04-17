package utils

import (
	"bytes"
	"encoding/hex"
	"io"
	"fmt"

	crypto "github.com/ProtonMail/go-crypto/openpgp"
)

var (
	pgpEntity *crypto.Entity
	hmacKey   = []byte("super‑secret‑hmac‑key")
)

func init() {
	var err error
	pgpEntity, err = crypto.NewEntity("bank", "cards", "bank@example.com", nil)
	if err != nil {
		panic(err)
	}
}

func EncryptPGP(data string) (cipherHex, macHex string, err error) {
	buf := new(bytes.Buffer)
	w, _ := crypto.Encrypt(buf, []*crypto.Entity{pgpEntity}, nil, nil, nil)
	io.WriteString(w, data)
	w.Close()

	cipherHex = hex.EncodeToString(buf.Bytes())
	macHex = ComputeHMAC(cipherHex, hmacKey)
	return
}

func DecryptPGP(cipherHex, macHex string) (string, error) {
	if macHex != ComputeHMAC(cipherHex, hmacKey) {
		return "", fmt.Errorf("HMAC mismatch")
	}
	cipherBytes, _ := hex.DecodeString(cipherHex)
	md, err := crypto.ReadMessage(bytes.NewReader(cipherBytes), crypto.EntityList{pgpEntity}, nil, nil)
	if err != nil {
		return "", err
	}
	plain, _ := io.ReadAll(md.UnverifiedBody)
	return string(plain), nil
}
