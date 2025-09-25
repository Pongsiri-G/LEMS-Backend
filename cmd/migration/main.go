package main

import (
	"fmt"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/database"
	"github.com/rs/zerolog/log"
)

func main() {
	config := configs.NewConfig()
	db := database.NewPostgrest(config)

	err := db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		log.Fatal().Msgf("Migration failed: %v", err)
	}

	fmt.Println("Migration completed ✅")
}
