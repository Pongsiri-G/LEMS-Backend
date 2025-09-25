package di

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers"
	authHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/context"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/database"
	userRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	authSv "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth"
	"github.com/google/wire"
)

var ConfigSet = wire.NewSet(
	configs.NewConfig,
)

var InfrastructureSet = wire.NewSet(
	context.NewContext,
	database.NewPostgrest,
)

var RepositorySet = wire.NewSet(
	userRepo.NewUserRepository,
)

var ServiceSet = wire.NewSet(
	authSv.NewAuthService,
)

var HandlerSet = wire.NewSet(
	handlers.NewHandlers,
	authHd.NewAuthHandler,
)
