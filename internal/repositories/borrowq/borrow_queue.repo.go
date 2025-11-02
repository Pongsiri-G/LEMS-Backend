package borrowq

import (
	"context"
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type BorrowQueueRepository interface {
	Enqueue(ctx context.Context, q *models.BorrowQueue) error
	PeekOldest(ctx context.Context, itemID string) (*models.BorrowQueue, error)
	PopFront(ctx context.Context, itemID string) (*models.BorrowQueue, error)
	GetFront(ctx context.Context, itemID string) (*models.BorrowQueue, error)
	Dequeue(ctx context.Context, queueID uuid.UUID) error
	Count(ctx context.Context, itemID string) (int, error)
	GetOneByUserAndItem(ctx context.Context, itemID, userID string) (*models.BorrowQueue, error)
	GetQueueByID(ctx context.Context, queueID uuid.UUID, where interface{}, arg ...interface{}) (*models.BorrowQueue, error)
	EditQueue(ctx context.Context, q *models.BorrowQueue) error
	DeleteQueue(cts context.Context, queueID uuid.UUID) error
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
	queue, err := b.GetQueueByID(ctx, queueID, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to get queue by id")
		return err
	}
	// err = b.db.Where("queue_id", queueID).Delete(&models.BorrowQueue{}).Error
	queue.IsBorrow = true
	now := time.Now()
	queue.BorrowedAt = &now

	err = b.EditQueue(ctx, queue)
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
		Where("i.item_id = (?) AND bq.is_borrow = (?)", itemID, false).
		Order("bq.created_at ASC").
		Limit(1).
		Scan(&queue)

	if res.Error != nil {
		log.Error().Err(res.Error).Msg("failed to get borrow queue")

		switch res.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, res.Error
		}
	}

	return queue, nil
}

// PopFront implements BorrowQueueRepository.
func (b *borrowQueueRepository) PopFront(ctx context.Context, itemID string) (*models.BorrowQueue, error) {
	queue, err := b.GetFront(ctx, itemID)
	if err != nil {
		log.Error().Err(err).Msg("failed to pop front borrow queue")
		return nil, err
	}

	err = b.Dequeue(ctx, queue.QueueID)
	if err != nil {
		log.Error().Err(err).Msg("failed to dequeue borrow queue")
		return nil, err
	}

	return queue, nil
}

// GetFront implements BorrowQueueRepository.
func (b *borrowQueueRepository) GetFront(ctx context.Context, itemID string) (*models.BorrowQueue, error) {
	var queue *models.BorrowQueue
	res := b.db.WithContext(ctx).Table("borrow_queues AS bq").
		Select("bq.*").
		Joins("JOIN items i ON i.item_id = bq.item_id").
		Where("i.item_id = (?) AND bq.is_borrow = (?)", itemID, false).
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
func (b *borrowQueueRepository) GetOneByUserAndItem(ctx context.Context, itemID, userID string) (*models.BorrowQueue, error) {
	var queue *models.BorrowQueue
	res := b.db.WithContext(ctx).Table("borrow_queues AS bq").
		Select("bq.*").
		Where("bq.item_id = (?) AND bq.user_id = (?)", itemID, userID).
		Limit(1).
		Scan(&queue)

	if res.Error != nil {
		log.Error().Err(res.Error).Msg("failed to get borrow queue")

		switch res.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, res.Error
		}
	}

	return queue, nil
}

// GetMemberByUserAndQID implements BorrowQueueRepository.
func (b *borrowQueueRepository) GetOneByUserAndQID(ctx context.Context, queueID, userID string) (*models.BorrowQueue, error) {
	var queue *models.BorrowQueue
	res := b.db.WithContext(ctx).Table("borrow_queues AS bq").
		Select("bq.*").
		Where("bq.queue_id = (?) AND bq.user_id = (?)", queueID, userID).
		Scan(&queue)

	if res.Error != nil {
		log.Error().Err(res.Error).Msg("failed to get borrow queue")

		switch res.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, res.Error
		}
	}
	return queue, nil
}

// GetQueueByID implements BorrowQueueRepository.
func (b *borrowQueueRepository) GetQueueByID(ctx context.Context, queueID uuid.UUID, where interface{}, arg ...interface{}) (*models.BorrowQueue, error) {
	var queue *models.BorrowQueue
	res := b.db.WithContext(ctx).Where("queue_id = ?", queueID).Where(where, arg...).First(&queue)

	if res.Error != nil {
		log.Error().Err(res.Error).Msg("failed to get queue by id")
		
		switch res.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, res.Error
		}
	}

	return queue, nil
}

// EditQueue implements BorrowQueueRepository.
func (b *borrowQueueRepository) EditQueue(ctx context.Context, q *models.BorrowQueue) error {
	err := b.db.WithContext(ctx).Save(q).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to edit queue")
	}
	return err
}

// DeleteQueue implements BorrowQueueRepository.
func (b *borrowQueueRepository) DeleteQueue(ctx context.Context, queueID uuid.UUID) error {
	return  b.db.WithContext(ctx).Where("queue_id", queueID).Delete(&models.BorrowQueue{}).Error	
}
