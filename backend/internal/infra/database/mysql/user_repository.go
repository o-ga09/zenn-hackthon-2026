package mysql

import (
	"context"

	"gorm.io/gorm"

	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
)

type UserRepository struct{}

// Create 新規ユーザーを作成
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	if err := Ctx.GetDB(ctx).Create(user).Error; err != nil {
		return errors.Wrap(ctx, err)
	}
	return nil
}

// FindByID IDでユーザーを検索
func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := Ctx.GetDB(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(ctx, err)
		}
		return nil, errors.Wrap(ctx, err)
	}
	return &user, nil
}

// FindByUID Firebase UIDでユーザーを検索
func (r *UserRepository) FindByUID(ctx context.Context, uid string) (*domain.User, error) {
	var user domain.User
	if err := Ctx.GetDB(ctx).Where("uid = ?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(ctx, err)
		}
		return nil, errors.Wrap(ctx, err)
	}
	return &user, nil
}

// FindAll 全ユーザーを取得
func (r *UserRepository) FindAll(ctx context.Context, opts *domain.FindOptions) ([]*domain.User, error) {
	var users []*domain.User
	query := Ctx.GetDB(ctx)

	if opts.Limit <= 0 {
		opts.Limit = 100
	}

	if err := query.Limit(opts.Limit).Offset(opts.Offset).Find(&users).Error; err != nil {
		return nil, errors.Wrap(ctx, err)
	}
	return users, nil
}

// Update ユーザー情報を更新
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	// 存在確認
	if _, err := r.FindByID(ctx, user.ID); err != nil {
		return errors.Wrap(ctx, err)
	}

	// 更新
	if err := Ctx.GetDB(ctx).Save(user).Error; err != nil {
		return errors.Wrap(ctx, err)
	}
	return nil
}

// Delete ユーザーを削除（論理削除）
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	// 存在確認
	if _, err := r.FindByID(ctx, id); err != nil {
		return err
	}

	// 論理削除
	if err := Ctx.GetDB(ctx).Delete(&domain.User{}, id).Error; err != nil {
		return errors.Wrap(ctx, err)
	}
	return nil
}
