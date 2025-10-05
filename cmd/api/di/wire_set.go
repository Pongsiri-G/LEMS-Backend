package di

import (
	"github.com/google/wire"

	// Configs
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/middlewares"

	// Handlers
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers"
	authHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"
	borrowHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/borrow"
	itemHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/item"
	minioHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/minio"
	userHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/user"

	// Infrastructure
	authInfra "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/context"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/database"
	minioInfra "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/minio"

	// Repositories
	borrowRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
	itemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	minioRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/minio"
	userRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"

	// Services
	authSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth/strategy"
	borrowSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/borrow"
	itemSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/jwt"
	minioSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/minio"
	userSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/user"
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
	itemRepo.NewItemRepository,
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
	itemSvc.NewItemService,
)

// ---- Handlers ----

var HandlerSet = wire.NewSet(
	handlers.NewHandlers,
	authHd.NewAuthHandler,
	userHd.NewUserHandler,
	minioHd.NewFileHandler,
	borrowHd.NewBorrowHandler,
	itemHd.NewItemHandler,
)

// ---- Middlewares ----
var MiddlewareSet = wire.NewSet(
	middlewares.NewAuthMiddleware,
)
