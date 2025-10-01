package di

import (
	"github.com/google/wire"

	// Configs
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"

	// Handlers
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers"
	authHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"

	// Infrastructure
	authInfra "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/context"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/database"

	// Repositories
	userRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"

	// Services
	authSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth/strategy"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/jwt"
)

var ConfigSet = wire.NewSet(
	configs.NewConfig,
)

var InfrastructureSet = wire.NewSet(
	context.NewContext,
	database.NewPostgrest,
	authInfra.NewGoogleOAuthClient,
)

var RepositorySet = wire.NewSet(
	userRepo.NewUserRepository,
)

// ---- Strategies ----

var StrategySet = wire.NewSet(
	strategy.NewStrategyMap,
	strategy.NewLocalStrategy,
	strategy.NewGoogleStrategy,
)

// ---- Services ----

var ServiceSet = wire.NewSet(
	jwt.NewJWTService,
	authSvc.NewAuthService,
)

// ---- Handlers ----

var HandlerSet = wire.NewSet(
	handlers.NewHandlers,
	authHd.NewAuthHandler,
)
