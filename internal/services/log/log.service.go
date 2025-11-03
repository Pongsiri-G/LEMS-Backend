package log

import (
	"context"
	"fmt"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	LogRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/log"
	UserRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/rs/zerolog/log"
)

type Service interface {
	GetAll(ctx context.Context) ([]responses.AdminLogResponse, error)
}

type service struct {
	logRepo  LogRepo.Repository
	userRepo UserRepo.Repository
}

func NewLogService(
	logRepo LogRepo.Repository,
	userRepo UserRepo.Repository,
) Service {
	return &service{
		logRepo:  logRepo,
		userRepo: userRepo,
	}
}

// GetAll implements Service.
func (s *service) GetAll(ctx context.Context) ([]responses.AdminLogResponse, error) {
	logs, err := s.logRepo.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all logs")
		return nil, err
	}
	var result []responses.AdminLogResponse
	for _, logModel := range logs {
		var res responses.AdminLogResponse
		user, err := s.userRepo.FindByID(ctx, logModel.UserID.String())
		if err != nil {
			log.Warn().Err(err).Str("user_id", logModel.UserID.String()).Msg("user not found for log, using deleted user placeholder")
			res.UserName = "[Deleted User]"
		} else {
			res.UserName = user.UserFullName
		}

		res.LogID = logModel.LogID.String()
		res.UserID = logModel.UserID.String()
		res.LogType = logModel.LogType
		res.Created = utils.ToStringDateTime(logModel.CreatedAt)

		// Format message based on log type
		if logModel.LogMessage != nil {
			message := *logModel.LogMessage
			// For admin actions, message contains target user ID - format as "username: user_id"
			if s.isAdminAction(logModel.LogType) {
				targetUser, err := s.userRepo.FindByID(ctx, message)
				if err != nil {
					log.Warn().Err(err).Str("target_user_id", message).Msg("target user not found for admin action log")
					res.Message = fmt.Sprintf("[Deleted User]: %s", message)
				} else {
					res.Message = fmt.Sprintf("%s: %s", targetUser.UserFullName, message)
				}
			} else {
				res.Message = message
			}
		}

		result = append(result, res)
	}
	return result, nil
}

func (s *service) isAdminAction(logType enums.LogType) bool {
	return logType == enums.LogTypeAccept ||
		logType == enums.LogTypeReject ||
		logType == enums.LogTypeActivate ||
		logType == enums.LogTypeDeactivate ||
		logType == enums.LogTypeDelete ||
		logType == enums.LogTypeGrantAdmin ||
		logType == enums.LogTypeRevokeAdmin
}
