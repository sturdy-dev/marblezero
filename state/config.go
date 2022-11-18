package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
)

type Config struct {
	Name string `json:"name"`

	// self saveable
	storagePath StoragePath `json:"-"`
}

type StoragePath string

func NewStoragePath() (StoragePath, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find config dir: %w", err)
	}

	configDir := path.Join(homeDir, ".config", "marblezero")
	if err := os.MkdirAll(configDir, 0777); err != nil {
		return "", fmt.Errorf("failed to create ~/.config/marblezero directory: %w", err)
	}

	return StoragePath(configDir), nil
}

func LoadConfig(storagePath StoragePath) (*Config, error) {
	contents, err := os.ReadFile(path.Join(string(storagePath), "config.json"))
	if errors.Is(err, os.ErrNotExist) {
		return &Config{
			storagePath: storagePath,
		}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(contents, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cfg.storagePath = storagePath

	return &cfg, nil
}

func (c *Config) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path.Join(string(c.storagePath), "config.json"), data, 0660); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	return nil
}
