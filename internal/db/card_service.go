package db

import (
	"banking-api/models"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"time"
)

func CreateCardService(db *sql.DB, userID int, accountID int) (*models.Card, error) {
	cardNumber, err := generateCardNumber()
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации номера карты: %v", err)
	}

	cvv, err := generateCVV()
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации CVV: %v", err)
	}

	expiration := time.Now().AddDate(3, 0, 0) // тип time.Time

	var cardID int
	createdAt := time.Now()
	err = db.QueryRow(`
		INSERT INTO cards (account_id, card_number, expiration_date, cvv, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, accountID, cardNumber, expiration, cvv, createdAt).Scan(&cardID)

	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении карты в базу: %v", err)
	}

	return &models.Card{
		ID:             cardID,
		AccountID:      accountID,
		CardNumber:     cardNumber,
		ExpirationDate: expiration,
		CVV:            cvv,
		CreatedAt:      createdAt,
	}, nil
}

func generateCardNumber() (string, error) {
	base := "400000"
	for len(base) < 15 {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		base += n.String()
	}

	sum := 0
	for i := 0; i < len(base); i++ {
		digit := int(base[len(base)-1-i] - '0')
		if i%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	checkDigit := (10 - (sum % 10)) % 10

	return base + fmt.Sprint(checkDigit), nil
}

func generateCVV() (string, error) {
	bytes := make([]byte, 2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%03d", int(bytes[0])%1000), nil
}
