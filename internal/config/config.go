package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DB_Url            string `json:"db_url"`
	Current_User_Name string `json:"current_user_name"`
}

func Read() (Config, error) {
	jsonFile, err := os.ReadFile(configFileName)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		return Config{}, err
	}

}

func (c *Config) SetUser() {

}

func getConfigFilePath() (string, error) {

}

func write(cfg Config) error {

}
