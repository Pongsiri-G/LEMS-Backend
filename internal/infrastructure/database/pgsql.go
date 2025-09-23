package database

import (
	"fmt"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres(cfg *configs.Config) {
	connection := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		cfg.DatabaseHost, 
		cfg.DatabaseUsername, 
		cfg.DatabasePassword, 
		cfg.DatabaseName,
		cfg.DatabasePort,
	)

	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to DB: %v", err)
	}

	DB = db
}