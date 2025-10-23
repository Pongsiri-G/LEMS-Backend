package user

import (
	"context"
	"errors"
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, userID string) (*models.User, error)
	Create(ctx context.Context, u *models.User) error
	FindById(ctx context.Context, id uuid.UUID) (*models.User, error)
	// ใช้สำหรับ provider-based login โดยไม่แยกตาราง
	// ข้อกำหนด: ลิงก์ด้วยอีเมลเสมอ ถ้าอีเมลมีอยู่แล้ว ให้ใช้งาน user เดิมและอัปเดต AuthProvider ตามความเหมาะสม
	FindOrCreateByProvider(ctx context.Context, provider enums.AuthProvider, email string, seed *models.User) (*models.User, error)
	UpdateStatus(ctx context.Context, userID uuid.UUID, status enums.UserStatus) error
	UpdateRole(ctx context.Context, userID uuid.UUID, role enums.UserRole) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	SoftDelete(ctx context.Context, userID uuid.UUID) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).
		Where("user_email = ?", email).
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).First(&u, "user_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) Create(ctx context.Context, u *models.User) error {
	now := time.Now()
	u.UserID = uuid.New()
	u.CreatedAt = now
	u.UpdatedAt = now
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *repository) FindOrCreateByProvider(ctx context.Context, provider enums.AuthProvider, email string, seed *models.User) (*models.User, error) {
	var u models.User
	tx := r.db.WithContext(ctx)

	err := tx.Where("user_email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if seed == nil {
				seed = &models.User{}
			}
			seed.UserID = uuid.New()
			seed.UserEmail = email
			seed.AuthProvider = provider
			now := time.Now()
			seed.CreatedAt = now
			seed.UpdatedAt = now
			if err := tx.Create(seed).Error; err != nil {
				return nil, err
			}
			return seed, nil
		}
		return nil, err
	}
	if u.AuthProvider != provider {
		u.AuthProvider = provider
		u.UpdatedAt = time.Now()
		if err := tx.Save(u).Error; err != nil {
			return nil, err
		}
	}
	return &u, nil
}

func (r *repository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("last_logged_in", now).Error
}

func (r *repository) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "user_id = ?", userID).Error
}

func (r *repository) UpdateStatus(ctx context.Context, userID uuid.UUID, status enums.UserStatus) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("user_id = ?", userID).Update("user_status", status).Error
}

func (r *repository) UpdateRole(ctx context.Context, userID uuid.UUID, role enums.UserRole) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("user_id = ?", userID).Update("user_role", role).Error
}

func (r *repository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	q := r.db.WithContext(ctx).Find(&users)
	if q.Error != nil {
		return nil, q.Error
	}
	return users, nil
// FindByID implements Repository.
func (r *repository) FindByID(ctx context.Context, userID string) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
