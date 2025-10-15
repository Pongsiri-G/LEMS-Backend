package main

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/cmd/api/di"
	"github.com/rs/zerolog/log"
)

func main() {
	server, err := di.InitializeAPI()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to initialize API")
	}
	if err := server.Start(); err != nil {
		log.Panic().Err(err).Msg("Failed to start server")
	}
}
