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
	JWT          JWT
	Database     DBConfig     `envPrefix:"DATABASE_"`
	Google       GoogleConfig `envPrefix:"GOOGLE_"`
	PG           PGConfig     `envPrefix:"PG_"`
	MINIO        MINIOConfig  `envPrefix:"MINIO_"`
	Email        Email
}

type DBConfig struct {
	Host     string `env:"HOST"`
	Name     string `env:"NAME"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Port     string `env:"PORT"`
}

type GoogleConfig struct {
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectURL  string `env:"REDIRECT_URL"`
}

type PGConfig struct {
	Email    string `env:"EMAIL"`
	Password string `env:"PASSWORD"`
}

type MINIOConfig struct {
	Endpoint string `env:"ENDPOINT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Bucket   string `env:"BUCKET"`
	UseSSL   bool   `env:"USE_SSL"`
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
