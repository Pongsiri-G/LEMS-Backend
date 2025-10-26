package database

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories"
	"gorm.io/gorm"
)

// gormTxKey is a context key used to detect nested transactions.
type gormTxKey struct{}

type gormTransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) repositories.TransactionManager {
	return &gormTransactionManager{db: db}
}

func (tm *gormTransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if tx := ctx.Value(gormTxKey{}); tx != nil {
		return fn(ctx)
	}

	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxWithTx := context.WithValue(ctx, gormTxKey{}, tx)
		return fn(ctxWithTx)
	})
}

func FromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(gormTxKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return fallback
}

