package configs

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port 				int16 	`env:"PORT"`

	DatabaseHost		string	`env:"DATABASE_HOST"`
	DatabaseName     	string	`env:"DATABASE_NAME"`
	DatabaseUsername 	string	`env:"DATABASE_USERNAME"`
	DatabasePassword 	string	`env:"DATABASE_PASSWORD"`
	DatabasePort     	int16	`env:"DATABASE_PORT"`
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found or error loading it.")
	}

	config := &Config{}

	if err := env.Parse(config); err != nil {
		log.Panic().Err(err).Msg("Failed to parse environment variables")
	}

	return nil
}
