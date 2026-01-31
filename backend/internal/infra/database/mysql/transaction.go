package mysql

import (
	"context"

	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"gorm.io/gorm"
)

type TransactionManager struct{}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{}
}

func (m *TransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return Ctx.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(Ctx.SetDB(ctx, tx))
	})
}
