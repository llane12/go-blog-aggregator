package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (cfg Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(cfg)
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(homeDir, configFileName)
	return p, nil
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonData, 0666)
}
