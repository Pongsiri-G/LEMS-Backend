package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserFilter struct {
	Status *enums.UserStatus
	Role   *enums.UserRole
	Search *string
	SortBy string
}

type Repository interface {
	// --- Query ---
	FindById(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	List(ctx context.Context, filter UserFilter) ([]models.User, error)

	// --- Command ---
	Create(ctx context.Context, u *models.User) error
	// ใช้สำหรับ provider-based login โดยไม่แยกตาราง
	FindOrCreateByProvider(ctx context.Context, provider enums.AuthProvider, email string, seed *models.User) (*models.User, error)
	UpdateStatus(ctx context.Context, userID uuid.UUID, status enums.UserStatus) error
	UpdateRole(ctx context.Context, userID uuid.UUID, role enums.UserRole) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	SoftDelete(ctx context.Context, userID uuid.UUID) error

	// --- Regacy --
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).
		Where("LOWER(user_email) = ?", email).
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *repository) List(ctx context.Context, f UserFilter) ([]models.User, error) {
	q := r.db.WithContext(ctx).Model(&models.User{})
	if f.Status != nil {
		q = q.Where("user_status = ?", *f.Status)
	}
	if f.Role != nil {
		q = q.Where("user_role = ?", *f.Role)
	}
	if f.Search != nil && *f.Search != "" {
		like := "%" + strings.ToLower(*f.Search) + "%"
		q = q.Where("LOWER(user_email) LIKE ? OR LOWER(user_full_name) LIKE ?", like, like)
	}
	if f.SortBy != "" {
		q = q.Order(f.SortBy)
	}
	var users []models.User
	if err := q.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
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
	res := r.db.WithContext(ctx).Delete(&models.User{}, "user_id = ?", userID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) UpdateStatus(ctx context.Context, userID uuid.UUID, status enums.UserStatus) error {
	res := r.db.WithContext(ctx).Model(&models.User{}).Where("user_id = ?", userID).Update("user_status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) UpdateRole(ctx context.Context, userID uuid.UUID, role enums.UserRole) error {
	res := r.db.WithContext(ctx).Model(&models.User{}).Where("user_id = ?", userID).Update("user_role", role)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	q := r.db.WithContext(ctx).Find(&users)
	if q.Error != nil {
		return nil, q.Error
	}
	return users, nil
}
