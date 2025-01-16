package config

import (
	"encoding/json"
	"log"
	"os"
)

const CONFIG_FILE_NAME = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	err := write(cfg)
	if err != nil {
		return err
	}
	return nil
}

func Read() Config {
	file_path, error := getConfigFilePath()
	if error != nil {
		log.Fatal("Error trying to get Config file path", error)
	}
	content, error := os.ReadFile(file_path)
	if error != nil {
		log.Fatal("Error while trying to read .gatorconfig.json at HOME, ", error)
	}

	var config Config
	error = json.Unmarshal(content, &config)
	if error != nil {
		log.Fatal("Error while trying to decode data", error)
	}

	return config
}

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error trying to read HOME dir", err)
		return home_dir, nil
	}

	file_path := home_dir + "/" + CONFIG_FILE_NAME

	return file_path, nil
}

func write(cfg *Config) error {
	file_path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(file_path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
