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
	borrowqHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/borrowq"
	itemHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/item"
	logHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/log"
	minioHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/minio"
	requestHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/request"
	tagHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/tag"
	userHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/user"
	wsHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/ws"

	// Infrastructure
	authInfra "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/context"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/database"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/email"
	minioInfra "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/minio"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/ws"

	// Repositories
	borrowRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
	borrowqRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrowq"
	itemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	itemRequestedRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_requested"
	itemsetRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_set"

	// logRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/log"
	logsystem "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/log"
	minioRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/minio"
	requestRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/request"
	tagRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/tag"
	userRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"

	// Services
	adminSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/admin"
	authSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth/strategy"
	borrowSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/borrow"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/borrowq"
	itemSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item"
	logSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/log"
	minioSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/minio"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/noti"
	requestSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/request"
	tagSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/tag"
	userSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/user"
)

var ConfigSet = wire.NewSet(
	configs.NewConfig,
)

var InfrastructureSet = wire.NewSet(
	ws.NewHub,
	context.NewContext,
	database.NewTransactionManager,
	database.NewPostgrest,
	minioInfra.NewMinioConnection,
	authInfra.NewGoogleOAuthClient,
	email.NewSMTPGoogle,
)

var RepositorySet = wire.NewSet(
	userRepo.NewUserRepository,
	minioRepo.NewMinioRepository,
	borrowRepo.NewBorrowLogRepository,
	itemRepo.NewItemRepository,
	tagRepo.NewTagRepository,
	itemsetRepo.NewItemSetRepository,
	logsystem.NewLogRepository,
	requestRepo.NewRepository,
	itemRequestedRepo.NewItemRequestedRepository,
	borrowqRepo.NewBorrowQueueRepository,
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
	noti.NewNotificationSubject,
	noti.ProvideSubjectWithObservers,
	itemSvc.NewItemService,
	tagSvc.NewTagService,
	requestSvc.NewRequestService,
	borrowq.NewBorrowQueueService,
	logSvc.NewLogService,
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
	requestHd.NewRequestHandler,
	borrowqHd.NewBorrowQueueHandler,
	wsHd.NewWsHandler,
	logHd.NewLogHandler,
)

// ---- Middlewares ----
var MiddlewareSet = wire.NewSet(
	middlewares.NewAuthMiddleware,
	middlewares.NewRbacMiddleware,
)
