package config

import (
	"fmt"
	"io"
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/" + configFileName, nil
}

func write(cfg Config) error {
	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent Error: %v", err)
	}

	configFilePathName, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonFile, err := os.Create(configFilePathName)
	if err != nil {
		return fmt.Errorf("os.Create Error: %v", err)
	}
	defer jsonFile.Close()

	if _, err := jsonFile.Write(jsonData); err != nil {
		return fmt.Errorf("Error Writing File: %v", err)
	}

	return nil
}
func Read() (Config, error) {
	configFilePathName, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(configFilePathName)
	if err != nil {
		return Config{}, fmt.Errorf("os.Open Error: %v", err)
	}
	defer jsonFile.Close()


	fileBytes,  err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, fmt.Errorf("io.ReadAll Error: %v", err)
	}

	var config Config
	if err := json.Unmarshal(fileBytes, &config); err != nil {
		return Config{}, fmt.Errorf("json.Unmarshal Error: %v", err)
	}

	return config, nil
}

func (c *Config) SetUser(username string) error {
	if len(username) < 1 {
		return fmt.Errorf("Please Provide a Username")
	}

	c.CurrentUserName = username
	if err := write(*c); err != nil {
		return err
	}

	return nil
}
