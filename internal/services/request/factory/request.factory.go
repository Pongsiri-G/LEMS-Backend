package factory

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
)

type Requestable interface {
	CreateRequest(ctx context.Context, req requests.CreateRequest) error
	EditRequest(ctx context.Context, req requests.EditRequest) error
}
