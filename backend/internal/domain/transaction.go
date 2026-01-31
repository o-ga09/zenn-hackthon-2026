package domain

import "context"

// ITransactionManager はトランザクションを管理するインターフェース
type ITransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
