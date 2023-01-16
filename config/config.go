package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	IsProd            bool   `toml:"is_prod" env:"IS_PROD" envDefault:"false"`
	DSN               string `toml:"dsn" env:"DSN" envDefault:"postgres://shopteam:123@localhost:5432/shop?sslmode=disable"`
	Port              string `toml:"port" env:"PORT" envDefault:":8000"` //80 8000
	FsPath            string `toml:"fs_path" env:"FS_PATH" envDefault:"./static/files/"`
	ServerURL         string `toml:"server_url" env:"SERVER_URL" envDefault:"http://localhost:8000"` // cozydragon.online localhost:8000
	Timeout           int    `toml:"timeout" env:"TIMEOUT" envDefault:"5"`
	CashHost          string `toml:"cash_host" env:"CASH_HOST" envDefault:"localhost"`
	CashPort          string `toml:"cash_port" env:"CASH_PORT" envDefault:"6379"`
	CashTTL           int    `toml:"cash_ttl" env:"CASH_TTL" envDefault:"30"`
	LogLevel          string `toml:"log_level" env:"LOG_LEVEL" envDefault:"debug"`
	ReadTimeout       int    `toml:"read_timeout" env:"READ_TIMEOUT" envDefault:"30"`
	WriteTimeout      int    `toml:"write_timeout" env:"WRITE_TIMEOUT" envDefault:"30"`
	ReadHeaderTimeout int    `toml:"read_header_timeout" env:"READ_HEADER_TIMEOUT" envDefault:"30"`
}

// NewConfig() initializes the configuration
func NewConfig() (*Config, error) {
	var configPath string

	// The flag allows to specify the path to the folder with the configuration file in .toml format.
	flag.StringVar(&configPath, "config-path", "", "path to file in .toml format")
	flag.Parse()

	var cfg = Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("can't load environment variables: %s", err)
	}
	if configPath != "" {
		_, err := toml.DecodeFile(configPath, &cfg)
		if err != nil {
			log.Fatalf("can't load configuration file: %s", err)
		}
	}

	configBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Configuration:", string(configBytes))
	return &cfg, nil
}
