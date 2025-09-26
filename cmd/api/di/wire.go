//go:build wireinject
// +build wireinject

package di

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/cmd/api/server"
	"github.com/google/wire"
)

func InitializeAPI() *server.EchoServer {
	wire.Build(
		ConfigSet,
		InfrastructureSet,
		RepositorySet,
		StrategySet,
		ServiceSet,
		HandlerSet,
		server.NewEchoServer,
	)

	return &server.EchoServer{}
}
