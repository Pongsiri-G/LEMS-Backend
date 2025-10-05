package di

import (
	"github.com/google/wire"

	// Configs
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/middlewares"

	// Handlers
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers"
	authHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"
	userHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/user"
	borrowHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/borrow"
	minioHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/minio"

	// Infrastructure
	authInfra "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/context"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/database"
	minioInfra "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/minio"

	// Repositories
	borrowRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
	minioRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/minio"
	userRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"

	// Services
	authSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth/strategy"
	userSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/user"
	borrowSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/borrow"
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
	borrowRepo.NewBorrowLogRepository,
)

// ---- Strategies ----

var StrategySet = wire.NewSet(
	strategy.NewStrategyMap,
	strategy.NewLocalStrategy,
	strategy.NewGoogleStrategy,
)

// ---- Services ----

var ServiceSet = wire.NewSet(
	authSvc.NewAuthService,
	userSvc.NewUserService,
	minioSvc.NewMinioService,
	borrowSvc.NewBorrowService,
)

// ---- Handlers ----

var HandlerSet = wire.NewSet(
	handlers.NewHandlers,
	authHd.NewAuthHandler,
	userHd.NewUserHandler,
	minioHd.NewFileHandler,
	borrowHd.NewBorrowHandler,
)

// ---- Middlewares ----
var MiddlewareSet = wire.NewSet(
	middlewares.NewAuthMiddleware,
)
