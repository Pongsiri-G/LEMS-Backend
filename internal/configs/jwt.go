package configs

type JWT struct {
	JwtSecret              string `env:"JWT_SECRET"`
	JwtExpirationMinutes   string `env:"JWT_EXPIRATION_MINUTES"`
	RefreshSecret          string `env:"REFRESH_SECRET"`
	RefreshExpirationHours string `env:"REFRESH_EXPIRATION_HOURS"`
}
