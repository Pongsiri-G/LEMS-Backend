package factory

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	"github.com/google/uuid"
)

type ItemFactory interface {
	CreateItem(ctx context.Context) (*models.Item, error)
}

type ItemResponseFactory interface {
	ToResponse(ctx context.Context, item *models.Item, children *[]models.Item) (*responses.ItemResponse, error)
}

type Borrowable interface {
	BorrowItem(ctx context.Context, userID uuid.UUID, item *models.Item, children *[]models.Item) error
	ReturnItem(ctx context.Context, borrowLog *models.BorrowLog, children *[]models.BorrowLog) error
}
