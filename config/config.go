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
		DB
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
		DSN      string
		User     string `env-required:"true" env:"DB_USER"`
		Password string `env-required:"true" env:"DB_PASSWORD"`
		Host     string `env-required:"true" env:"DB_HOST"`
		Port     string `env-required:"true" env:"DB_PORT"`
		Name     string `env-required:"true" env:"DB_NAME"`
	}
)

// Returns app config loaded from YAML-file.
func New() (*Config, error) {
	cfg := &Config{}

	// read ENV variables
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read env-variables: %w", err)
	}
	// read YAML config file
	if err := cleanenv.ReadConfig("./config.yml", cfg); err != nil {
		return nil, fmt.Errorf("read yaml config file: %w", err)
	}
	cfg.DB.DSN = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&timeout=10s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)
	return cfg, nil
}
