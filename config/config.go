// Package config provides loading config data from external sources
// like env variables, yaml-files etc.
package config

import "time"

type (
	Config struct {
		App      `yaml:"app"`
		Hardware `yaml:"hardware"`
		Other    `yaml:"other"`
		DB
	}

	App struct {
		LogLevel         string `yaml:"log_level"`
		LogFormat        string `yaml:"log_format"`
		GreetingsImgPath string `yaml:"greetings_img_path"`
	}

	Hardware struct {
		Oled    `yaml:"oled"`
		Buttons `yaml:"buttons"`
	}

	Oled struct {
		Bus             string        `yaml:"bus"`
		Width           int           `yaml:"width"`
		Height          int           `yaml:"height"`
		UpdatesDuration time.Duration `yaml:"updates_duration"`
	}

	Buttons struct {
		Up                string        `yaml:"up"`
		Down              string        `yaml:"down"`
		Escape            string        `yaml:"escape"`
		Enter             string        `yaml:"enter"`
		CheckAliveTimeout time.Duration `yaml:"check_alive_timeout"`
	}

	Other struct {
		Station `yaml:"station"`
		Sensors `yaml:"sensors"`
	}

	Station struct {
		ServiceName string `yaml:"service_name"`
		ConfigPath  string `yaml:"config_path"`
	}

	Sensors struct {
		DataUpdatesDuration time.Duration `yaml:"data_updates_duration"`
	}

	DB struct {
		DSN       string
		User      string `env-required:"true" env:"DB_USER"`
		Password  string `env-required:"true" env:"DB_PASSWORD"`
		Name      string `env-required:"true" env:"DB_NAME"`
		Host      string `env:"DB_HOST" env-default:"127.0.0.1"`
		Port      string `env:"DB_PORT" env-default:"3306"`
		Socket    string `env:"DB_SOCKET" env-default:"/var/run/mysqld/mysqld.sock" env-description:"unix socket"`
		UseSocket bool   `env:"DB_USE_SOCKET" env-default:"false" env-description:"use unix-socket instead of tcp connection (default: false)"`
	}
)

// Returns app config loaded from YAML-file.
func New() (*Config, error) {
	cfg := &Config{}

	// // read ENV variables
	// if err := cleanenv.ReadEnv(cfg); err != nil {
	// 	return nil, fmt.Errorf("read env-variables: %w", err)
	// }
	// // read YAML config file
	// if err := cleanenv.ReadConfig("./config.yml", cfg); err != nil {
	// 	return nil, fmt.Errorf("read yaml config file: %w", err)
	// }

	// var dbConnParams string
	// if cfg.DB.UseSocket {
	// 	dbConnParams = fmt.Sprintf("unix(%s)", cfg.DB.Socket)
	// } else {
	// 	dbConnParams = fmt.Sprintf("tcp(%s:%s)", cfg.DB.Host, cfg.DB.Port)
	// }
	// cfg.DB.DSN = fmt.Sprintf(
	// 	"%s:%s@%s/%s?parseTime=true&timeout=10s",
	// 	cfg.DB.User,
	// 	cfg.DB.Password,
	// 	dbConnParams,
	// 	cfg.DB.Name,
	// )
	return cfg, nil
}
