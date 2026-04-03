package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	homeDirectory, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	jsonFile, err := os.ReadFile(homeDirectory)
	var config Config
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	filePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath = filePath + "/" + configFileName
	return filePath, nil
}

func write(cfg Config) error {
	jsonContent, err := json.Marshal(cfg)
	homeDirectory, err := getConfigFilePath()
	err = os.WriteFile(homeDirectory, jsonContent, 0777)
	if err != nil {
		return err
	}

	return nil
}
