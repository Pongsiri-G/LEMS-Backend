package log

import (
	"context"

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
			log.Error().Err(err).Msg("failed to get user by id")
			return nil, err
		}
		res.LogID = logModel.LogID.String()
		res.UserID = logModel.UserID.String()
		res.LogType = logModel.LogType
		res.Created = utils.ToStringDateTime(logModel.CreatedAt)
		res.UserName = user.UserFullName
		if logModel.LogMessage != nil {
			res.Message = *logModel.LogMessage
		}

		result = append(result, res)
	}
	return result, nil
}
