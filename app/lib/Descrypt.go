package lib

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	lg "app.go/app/lib/logger"
)

func encrypt(bot_token string) string {
	// Calculate SHA256 hash
	hash := sha256.Sum256([]byte(bot_token))

	// Convert hash to string
	return hex.EncodeToString(hash[:])
}

func decrypt(hash string, log *lg.Logger) string {
	// Convert base64 string to bytes
	bytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		log.Fatal(err, "Ошибка декодирования пароля:")
		return ""
	}
	return string(bytes)
}

func DecryptBotToken(bot_token, bot_hash string, log *lg.Logger) string {
	hash_config := encrypt(bot_token)

	if hash_config == bot_hash {
		log.Info(fmt.Sprintf("Пароль совпал: %s", hash_config))
		return decrypt(bot_token, log)
	} else {
		log.Fatal(fmt.Errorf("неверное совпадение пароля"), "")
		return ""
	}
}
