package db

import (
	"banking-api/models"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var pgpPublicKey = []byte(`-----BEGIN PGP PUBLIC KEY BLOCK-----
... ваш публичный PGP-ключ ...
-----END PGP PUBLIC KEY BLOCK-----`)

func encryptPGP(plainText string) (string, error) {
	return plainText, nil
}

func hashCVV(cvv string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func computeHMAC(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func SaveCard(db *sql.DB, accountID int, cardNumber, cvv string, expirationDate time.Time) (*models.Card, error) {
	encryptedNumber, err := encryptPGP(cardNumber)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования номера карты: %v", err)
	}

	hashedCVV, err := hashCVV(cvv)
	if err != nil {
		return nil, fmt.Errorf("ошибка хеширования CVV: %v", err)
	}

	hmacValue := computeHMAC(cardNumber, []byte("card_hmac_secret"))

	var id int
	createdAt := time.Now()
	err = db.QueryRow(`
		INSERT INTO cards (account_id, card_number, expiration_date, cvv, hmac, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`, accountID, encryptedNumber, expirationDate, hashedCVV, hmacValue, createdAt).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения карты: %v", err)
	}

	return &models.Card{
		ID:             id,
		AccountID:      accountID,
		CardNumber:     encryptedNumber,
		ExpirationDate: expirationDate,
		CVV:            hashedCVV,
		HMAC:           hmacValue,
		CreatedAt:      createdAt,
	}, nil
}
