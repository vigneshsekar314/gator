package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	fullpath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}
	rawFileContent, err := os.ReadFile(fullpath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(rawFileContent, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	fullpath, err := getConfigPath()
	if err != nil {
		return err
	}
	confBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	os.WriteFile(fullpath, confBytes, 0666)
	return nil
}

func getConfigPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", homedir, configFileName), nil
}
