package lib

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"log"
)

func encrypt(bot_token string) string {
	// Calculate SHA256 hash
	hash := sha256.Sum256([]byte(bot_token))

	// Convert hash to string
	return hex.EncodeToString(hash[:])
}

func decrypt(hash string) string {
	// Convert base64 string to bytes
	bytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		log.Fatalf("Ошибка декодирования пароля:\n%v", err)
		return ""
	}
	return string(bytes)
}

func GetDecryptBotToken(bot_token, bot_hash string) string {
	hash_config := encrypt(bot_token)

	if hash_config == bot_hash {
		log.Printf("Пароль совпал: %s", hash_config)
		return decrypt(bot_token)
	} else {
		log.Fatal("Неверное совпадение пароля")
		return ""
	}
}
