package config

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

type Config struct {
    DbUrl string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
    configPath, err := getConfigPath()
    if err != nil {
        return Config{}, err
    }

    data, err := os.ReadFile(configPath)
    if err != nil {
        return Config{}, err
    }

    config := Config{}

    if err := json.Unmarshal(data, &config); err != nil {
        return Config{}, err
    }

    return config, nil
}

func getConfigPath() (string, error) {
    homePath, err := os.UserHomeDir()
    if err != nil {
        log.Printf("Failed to fetch home path: %v\n", err)
        return "", err
    }
    return path.Join(homePath, configFileName), nil
}

func write(cfg Config) error {
    configPath, err := getConfigPath()
    if err != nil {
        return err
    }

    data, err := json.Marshal(cfg)
    if err != nil {
        return err
    }

    return os.WriteFile(configPath, data, 0666)
}

func (c *Config) SetUser(username string) error {
    c.CurrentUserName = username

    err := write(*c)
    if err != nil {
        return err
    }

    return nil
}
