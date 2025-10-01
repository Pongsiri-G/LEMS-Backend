package di

import (
	"github.com/google/wire"

	// Configs
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"

	// Handlers
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers"
	authHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"
	minioHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/minio"

	// Infrastructure
	authInfra "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/context"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/database"
	minioInfra "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/minio"

	// Repositories
	minioRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/minio"
	userRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"

	// Services
	authSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth/strategy"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/jwt"
	minioSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/minio"
)

var ConfigSet = wire.NewSet(
	configs.NewConfig,
)

var InfrastructureSet = wire.NewSet(
	context.NewContext,
	database.NewPostgrest,
	minioInfra.NewMinioConnection,
	authInfra.NewGoogleOAuthClient,
)

var RepositorySet = wire.NewSet(
	userRepo.NewUserRepository,
	minioRepo.NewMinioRepository,
)

// ---- Strategies ----

type strategyDeps struct {
	Local  *strategy.LocalStrategy
	Google *strategy.GoogleStrategy
}

func newStrategyMap(d strategyDeps) map[string]strategy.AuthStrategy {
	return map[string]strategy.AuthStrategy{
		"local":  d.Local,
		"google": d.Google,
	}
}

var StrategySet = wire.NewSet(
	strategy.NewLocalStrategy,
	strategy.NewGoogleStrategy,
	wire.Struct(new(strategyDeps), "*"),
	newStrategyMap,
)

// ---- Services ----

var ServiceSet = wire.NewSet(
	jwt.NewJWTService,
	authSvc.NewAuthService,
	minioSvc.NewMinioService,
)

// ---- Handlers ----

var HandlerSet = wire.NewSet(
	handlers.NewHandlers,
	authHd.NewAuthHandler,
	minioHd.NewFileHandler,
)
