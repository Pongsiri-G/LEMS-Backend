package borrowq

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type BorrowQueueRepository interface {
	Enqueue(ctx context.Context, q *models.BorrowQueue) error
	PeekOldest(ctx context.Context, itemID string) (*models.BorrowQueue, error)
	Dequeue(ctx context.Context, queueID uuid.UUID) error
	Count(ctx context.Context, itemID string) (int, error)
	GetMemberByUserAndItem(ctx context.Context, itemID string, userID string) (*models.BorrowQueue, error)
}

type borrowQueueRepository struct {
	db *gorm.DB
}

func NewBorrowQueueRepository(db *gorm.DB) BorrowQueueRepository {
	return &borrowQueueRepository{
		db: db,
	}
}

// Count implements BorrowQueueRepository.
func (b *borrowQueueRepository) Count(ctx context.Context, itemID string) (int, error) {
	var num int
	result := b.db.WithContext(ctx).Table("borrow_queues AS bq").
		Select("COUNT(bq.*)").
		Joins("JOIN items i ON i.item_id = bq.item_id").
		Where("i.item_id = (?) bq.is_borrow = (?)", itemID, false).
		Scan(&num)

	if result.Error != nil {
		return 0, result.Error
	}

	return num, nil
}

// Dequeue implements BorrowQueueRepository.
func (b *borrowQueueRepository) Dequeue(ctx context.Context, queueID uuid.UUID) error {
	err := b.db.Where("queue_id", queueID).Delete(&models.BorrowQueue{}).Error

	if err != nil {
		log.Error().Err(err).Msg("failed to dequeue")
		return err
	}

	return nil
}

// Enqueue implements BorrowQueueRepository.
func (b *borrowQueueRepository) Enqueue(ctx context.Context, queue *models.BorrowQueue) error {
	tx := database.FromContext(ctx, b.db)
	queue.QueueID = uuid.New()
	return tx.Create(queue).Error
}

// PeekOldest implements BorrowQueueRepository.
func (b *borrowQueueRepository) PeekOldest(ctx context.Context, itemID string) (*models.BorrowQueue, error) {
	var queue *models.BorrowQueue
	res := b.db.WithContext(ctx).Table("borrow_queues AS bq").
		Select("bq.*").
		Joins("JOIN items i ON i.item_id = bq.item_id").
		Where("i.item_id = (?)", itemID).
		Order("bq.created_at ASC").
		Limit(1).
		Scan(&queue)

	if res.Error != nil {
		log.Error().Err(res.Error).Msg("failed to get borrow queue")

		return nil, res.Error
	}

	return queue, nil
}

// GetMemberByUser implements BorrowQueueRepository.
func (b *borrowQueueRepository) GetMemberByUserAndItem(ctx context.Context, itemID string, userID string) (*models.BorrowQueue, error) {
	var queue *models.BorrowQueue
	res := b.db.WithContext(ctx).Table("borrow_queues AS bq").
		Select("bq.*").
		Where("bq.item_id = (?) AND bq.user_id = (?)", itemID, userID).
		Limit(1).
		Scan(&queue)

	if res.Error != nil {
		log.Error().Err(res.Error).Msg("failed to get borrow queue")

		return nil, res.Error
	}

	return queue, nil
}
