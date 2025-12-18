package main

import (
	"errors"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	ApiKey  string `yaml:"api-key"`
	Webhook string `yaml:"webhook"`
}

func LoadConfig() (*Config, error) {
	file, err := os.OpenFile("config.yml", os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var cfg Config

	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.ApiKey == "" {
		return nil, errors.New("missing `api-key`")
	} else if cfg.Webhook == "" {
		return nil, errors.New("missing `webhook`")
	}

	return &cfg, nil
}
