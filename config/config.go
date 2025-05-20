// Package config provides loading config data from external sources
// like env variables, yaml-files etc.
package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App `yaml:"app"`
		DB  `yaml:"db"`
	}

	App struct {
		OledBus          string `yaml:"oled_bus"`
		GreetingsImgPath string `yaml:"greetings_img_path"`
	}

	DB struct {
	}
)

// Returns app config loaded from YAML-file.
func New() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig("./config.yml", cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
