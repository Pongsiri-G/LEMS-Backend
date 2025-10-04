package configs

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port         string   `env:"PORT"`
	AllowOrigins []string `env:"ALLOW_ORIGINS" envSeparator:","`

	DatabaseHost       string `env:"DATABASE_HOST"`
	DatabaseName       string `env:"DATABASE_NAME"`
	DatabaseUsername   string `env:"DATABASE_USERNAME"`
	DatabasePassword   string `env:"DATABASE_PASSWORD"`
	DatabasePort       string `env:"DATABASE_PORT"`
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `env:"GOOGLE_REDIRECT_URL"`
	JWT                JWT
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found or error loading it.")
	}

	config := &Config{}

	if err := env.Parse(config); err != nil {
		log.Panic().Err(err).Msg("Failed to parse environment variables")
	}

	fmt.Println(config)

	return config
}
