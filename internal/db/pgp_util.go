package db

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

func encryptPGP(plainText string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("не удалось получить домашнюю директорию: %v", err)
	}
	keyPath := filepath.Join(homeDir, "public_pgp_key.asc")
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения PGP-ключа из файла: %v", err)
	}

	block, err := armor.Decode(bytes.NewReader(keyData))
	if err != nil {
		return "", fmt.Errorf("ошибка декодирования PGP-ключа: %v", err)
	}

	entity, err := openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil {
		return "", fmt.Errorf("ошибка чтения сущности PGP: %v", err)
	}

	var buf bytes.Buffer
	armoredWriter, err := armor.Encode(&buf, "PGP MESSAGE", nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания armored writer: %v", err)
	}

	encryptWriter, err := openpgp.Encrypt(armoredWriter, []*openpgp.Entity{entity}, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка шифрования PGP: %v", err)
	}

	if _, err := io.WriteString(encryptWriter, plainText); err != nil {
		return "", fmt.Errorf("ошибка записи в шифратор: %v", err)
	}

	if err := encryptWriter.Close(); err != nil {
		return "", fmt.Errorf("ошибка закрытия шифратора: %v", err)
	}
	if err := armoredWriter.Close(); err != nil {
		return "", fmt.Errorf("ошибка закрытия armored writer: %v", err)
	}

	return buf.String(), nil
}
