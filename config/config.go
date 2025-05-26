// Package config provides loading config data from external sources
// like env variables, yaml-files etc.
package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"app"`
		Hardware `yaml:"hardware"`
		DB       `yaml:"db"`
	}

	App struct {
		GreetingsImgPath string `yaml:"greetings_img_path"`
	}

	Hardware struct {
		Oled    `yaml:"oled"`
		Buttons `yaml:"buttons"`
	}

	Oled struct {
		Bus string `yaml:"bus"`
	}

	Buttons struct {
		Up     string `yaml:"up"`
		Down   string `yaml:"down"`
		Escape string `yaml:"escape"`
		Enter  string `yaml:"enter"`
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
