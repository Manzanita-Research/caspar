package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".caspar.json"

type Config struct {
	URL         string `json:"url"`
	AdminAPIKey string `json:"admin_api_key"`
}

func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	return filepath.Join(home, configFileName), nil
}

func Load() (*Config, error) {
	// env vars override everything
	url := os.Getenv("CASPAR_URL")
	key := os.Getenv("CASPAR_ADMIN_API_KEY")
	if url != "" && key != "" {
		return &Config{URL: url, AdminAPIKey: key}, nil
	}

	path, err := Path()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not logged in — run `caspar auth login` first")
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// allow partial env var overrides
	if url != "" {
		cfg.URL = url
	}
	if key != "" {
		cfg.AdminAPIKey = key
	}

	if cfg.URL == "" || cfg.AdminAPIKey == "" {
		return nil, fmt.Errorf("incomplete config — run `caspar auth login`")
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := Path()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

func Delete() error {
	path, err := Path()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing config: %w", err)
	}

	return nil
}
