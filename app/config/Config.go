package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	lg "app.go/app/lib/logger"
)

type Config struct {
	Bot_secure struct {
		Bot_token string `yaml:"bot_token"`
		Bot_hash  string `yaml:"bot_hash"`
		Chat_id   int    `yaml:"chat_id"`
	} `yaml:"bot_secure"`
	Poll struct {
		Question   string   `yaml:"question"`
		AnswersYes []string `yaml:"answersYes"`
		AnswersNo  []string `yaml:"answersNo"`
	} `yaml:"poll"`
	Url          string `yaml:"url"`
	Path_to_pics string `yaml:"path_to_pics"`
}

func LoadConfig(log *lg.Logger) Config {
	// Create new struct
	var config Config

	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err, "Ошибка получения пути исполняемого файла:")
		return config
	}

	// Get executable dir
	execDir := filepath.Dir(execPath)

	// Set config path
	filename := filepath.Join(execDir, "config", "config.yaml")

	// Check if file exists
	_, err = os.Stat(filename)
	if err != nil {
		log.Fatal(err, fmt.Sprintf("Файл конфигурации '%s' не существует:", filename))
		return config
	}

	// Check if file is readable
	_, err = os.Open(filename)
	if err != nil {
		log.Fatal(err, fmt.Sprintf("Файл конфигурации '%s' недоступен для чтения:", filename))
		return config
	}

	// Read config file
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err, fmt.Sprintf("Ошибка чтения файла '%s':", filename))
		return config
	}

	// Parsing YAML to struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err, fmt.Sprintf("Ошибка парсинга файла '%s':", filename))
		return config
	}

	log.Info(fmt.Sprintf("Конфигурация успешно загружена из '%s'", filename))

	return config
}
