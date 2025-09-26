package user

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, u *models.User) error

	// ใช้สำหรับ provider-based login โดยไม่แยกตาราง
	// ข้อกำหนด: ลิงก์ด้วยอีเมลเสมอ ถ้าอีเมลมีอยู่แล้ว ให้ใช้งาน user เดิมและอัปเดต AuthProvider ตามความเหมาะสม
	FindOrCreateByProvider(ctx context.Context, provider enums.AuthProvider, email string, seed *models.User) (*models.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ---- Phase 1: dummy implementations ----

func (r *repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	// TODO: Phase 2 implement with GORM
	return nil, nil
}

func (r *repository) Create(ctx context.Context, u *models.User) error {
	// TODO: Phase 2 implement with GORM
	return nil
}

func (r *repository) FindOrCreateByProvider(ctx context.Context, provider enums.AuthProvider, email string, seed *models.User) (*models.User, error) {
	// TODO: Phase 2 implement with GORM
	return nil, nil
}
