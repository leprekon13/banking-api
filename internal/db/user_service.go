package db

import (
	"banking-api/models"
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

func RegisterUser(db *sql.DB, username, email, password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("ошибка хеширования пароля: %v", err)
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	}

	err = AddUser(db, user)
	if err != nil {
		return "", fmt.Errorf("ошибка регистрации пользователя: %v", err)
	}

	token, err := GenerateJWT(user)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации JWT: %v", err)
	}

	return token, nil
}

func LoginUserService(db *sql.DB, email, password string) (string, error) {
	user, err := LoginUser(db, email)
	if err != nil {
		return "", fmt.Errorf("ошибка входа: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("неверный пароль")
	}

	token, err := GenerateJWT(user)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации JWT: %v", err)
	}

	return token, nil
}

func GenerateJWT(user *models.User) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки .env файла: %v", err)
	}

	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
