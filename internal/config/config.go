package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("error getting user home directory: %w", err)
	}
	fullPath := filepath.Join(homeDir, configFileName)
	configBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config json: %w", err)
	}

	config := Config{}
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return Config{}, fmt.Errorf("error creating config struct: %w", err)
	}
	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home directory: %w", err)
	}

	configByte, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("error converting to json: %w", err)
	}

	fullPath := filepath.Join(homeDir, configFileName)
	if err := os.WriteFile(fullPath, configByte, 0644); err != nil {
		return fmt.Errorf("error writing config json to user home directory: %w", err)
	}

	return nil
}
