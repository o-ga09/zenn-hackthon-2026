package mysql

import (
	"context"

	"gorm.io/gorm"

	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
)

type NotificationRepository struct{}

// Create - 通知を作成
func (r *NotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	if err := Ctx.GetDB(ctx).Create(notification).Error; err != nil {
		return errors.Wrap(ctx, err)
	}
	return nil
}

// FindByID - IDで通知を取得
func (r *NotificationRepository) FindByID(ctx context.Context, id string) (*domain.Notification, error) {
	var notification domain.Notification
	userID := Ctx.GetCtxFromUser(ctx)
	if err := Ctx.GetDB(ctx).Where("id = ? AND user_id = ?", id, userID).First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(ctx, err)
		}
		return nil, errors.Wrap(ctx, err)
	}
	return &notification, nil
}

// FindByUserID - ユーザーIDで通知を取得（created_at降順）
func (r *NotificationRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	if err := Ctx.GetDB(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		return nil, errors.Wrap(ctx, err)
	}
	return notifications, nil
}

// MarkAsRead - 通知を既読にする
func (r *NotificationRepository) MarkAsRead(ctx context.Context, notification *domain.Notification) error {
	notification.Read = true
	if err := Ctx.GetDB(ctx).
		Updates(notification).Error; err != nil {
		return errors.Wrap(ctx, err)
	}
	return nil
}

// MarkAllAsRead - ユーザーの全通知を既読にする
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, notification *domain.Notification) (int64, error) {
	ctxWithSkip := Ctx.WithSkipOptimisticLock(ctx)

	result := Ctx.GetDB(ctxWithSkip).
		Where("user_id = ?", notification.UserID).
		Updates(&domain.Notification{Read: true})

	if result.Error != nil {
		return 0, errors.Wrap(ctx, result.Error)
	}

	return result.RowsAffected, nil
}

// Delete - 通知を削除
func (r *NotificationRepository) Delete(ctx context.Context, notification *domain.Notification) error {
	if err := Ctx.GetDB(ctx).Delete(notification).Error; err != nil {
		return errors.Wrap(ctx, err)
	}
	return nil
}

// CountUnread - 未読通知数を取得
func (r *NotificationRepository) CountUnread(ctx context.Context, userID string) (int64, error) {
	var count int64
	if err := Ctx.GetDB(ctx).Model(&domain.Notification{}).
		Where("user_id = ? AND `read` = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, errors.Wrap(ctx, err)
	}
	return count, nil
}

// DeleteAllByUserID - ユーザーの全通知を削除
func (r *NotificationRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	if err := Ctx.GetDB(ctx).
		Delete(&domain.Notification{}, &domain.Notification{UserID: userID}).Error; err != nil {
		return errors.Wrap(ctx, err)
	}

	return nil
}
