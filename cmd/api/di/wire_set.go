package di

import (
	"github.com/google/wire"

	// Configs
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/middlewares"

	// Handlers
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers"
	adminHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/admin"
	authHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"
	borrowHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/borrow"
	itemHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/item"
	minioHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/minio"
	tagHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/tag"
	userHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/user"

	// Infrastructure
	authInfra "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/context"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/database"
	minioInfra "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/minio"

	// Repositories
	borrowRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
	itemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	itemsetRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_set"
	logsystem "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/log"
	minioRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/minio"
	tagRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/tag"
	userRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"

	// Services
	adminSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/admin"
	authSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth/strategy"
	borrowSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/borrow"
	itemSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item"
	minioSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/minio"
	tagSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/tag"
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
	tagRepo.NewTagRepository,
	itemsetRepo.NewItemSetRepository,
	logsystem.NewLogRepository,
)

// ---- Strategies ----

var StrategySet = wire.NewSet(
	strategy.NewStrategyMap,
	strategy.NewLocalStrategy,
	strategy.NewGoogleStrategy,
)

// ---- Services ----

var ServiceSet = wire.NewSet(
	adminSvc.NewAdminService,
	authSvc.NewAuthService,
	userSvc.NewUserService,
	minioSvc.NewMinioService,
	borrowSvc.NewBorrowService,
	itemSvc.NewItemService,
	tagSvc.NewTagService,
)

// ---- Handlers ----

var HandlerSet = wire.NewSet(
	handlers.NewHandlers,
	adminHd.NewAdminHandler,
	authHd.NewAuthHandler,
	userHd.NewUserHandler,
	minioHd.NewFileHandler,
	borrowHd.NewBorrowHandler,
	itemHd.NewItemHandler,
	tagHd.NewTagHandler,
)

// ---- Middlewares ----
var MiddlewareSet = wire.NewSet(
	middlewares.NewAuthMiddleware,
	middlewares.NewRbacMiddleware,
)
