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

var hmacSecret = []byte("supersecretkey_for_hmac")

func SaveCard(db *sql.DB, accountID int, cardNumber, cvv string, expirationDate time.Time) (*models.Card, error) {
	hashedCVV, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("ошибка хеширования CVV: %v", err)
	}

	hmacValue := computeHMAC(cardNumber, hmacSecret)

	var id int
	createdAt := time.Now()
	err = db.QueryRow(`
		INSERT INTO cards (account_id, card_number, expiration_date, cvv, hmac, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`, accountID, cardNumber, expirationDate, string(hashedCVV), hmacValue, createdAt).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения карты: %v", err)
	}

	return &models.Card{
		ID:             id,
		AccountID:      accountID,
		CardNumber:     cardNumber,
		ExpirationDate: expirationDate,
		CVV:            "", // CVV не возвращается в ответе
		HMAC:           hmacValue,
		CreatedAt:      createdAt,
	}, nil
}

func computeHMAC(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func GetCardsByUserID(db *sql.DB, userID int) ([]models.Card, error) {
	rows, err := db.Query(`
		SELECT c.id, c.account_id, c.card_number, c.expiration_date, c.created_at, c.hmac
		FROM cards c
		JOIN accounts a ON c.account_id = a.id
		WHERE a.user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var card models.Card
		var hmac sql.NullString

		err := rows.Scan(
			&card.ID,
			&card.AccountID,
			&card.CardNumber,
			&card.ExpirationDate,
			&card.CreatedAt,
			&hmac,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %v", err)
		}

		if hmac.Valid {
			card.HMAC = hmac.String
		}

		cards = append(cards, card)
	}

	return cards, nil
}
