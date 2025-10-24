package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	SearchPaths []string `json:"search_paths"`
	MaxDepth    int      `json:"max_depth"`
	Theme       string   `json:"theme"`
}

func DefaultConfig() Config {
	return Config{
		SearchPaths: []string{
			"~/projects",
			"~/code",
			"~/work",
			"~/dev",
		},
		MaxDepth: 3,
		Theme:    "dark",
	}
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, ".config", "tsm")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

func Load() (Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return DefaultConfig(), err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return DefaultConfig(), err
	}

	return cfg, nil
}

func Save(cfg Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func SaveDefault() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		return nil
	}

	return Save(DefaultConfig())
}
