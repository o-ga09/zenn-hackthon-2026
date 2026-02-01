package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
)

type INotificationHandler interface {
	GetNotifications(ctx echo.Context) error
	MarkAsRead(ctx echo.Context) error
	MarkAllAsRead(ctx echo.Context) error
	DeleteNotification(ctx echo.Context) error
}

type NotificationHandler struct {
	notificationRepo domain.INotificationRepository
}

func NewNotificationHandler(notificationRepo domain.INotificationRepository) *NotificationHandler {
	return &NotificationHandler{
		notificationRepo: notificationRepo,
	}
}

// GetNotifications - 通知一覧取得
func (h *NotificationHandler) GetNotifications(c echo.Context) error {
	ctx := c.Request().Context()
	userID := context.GetCtxFromUser(ctx)

	notifications, err := h.notificationRepo.FindByUserID(ctx, userID)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	unreadCount, err := h.notificationRepo.CountUnread(ctx, userID)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"notifications": notifications,
		"unread_count":  unreadCount,
	})
}

// MarkAsRead - 通知を既読にする
func (h *NotificationHandler) MarkAsRead(c echo.Context) error {
	ctx := c.Request().Context()
	userID := context.GetCtxFromUser(ctx)
	notificationID := c.Param("id")

	// 通知の所有権を確認
	notification, err := h.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Notification not found",
		})
	}

	if notification.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Access denied",
		})
	}

	if err := h.notificationRepo.MarkAsRead(ctx, notificationID); err != nil {
		return errors.Wrap(ctx, err)
	}

	// 更新後の通知を返す
	updatedNotification, err := h.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.JSON(http.StatusOK, updatedNotification)
}

// MarkAllAsRead - 全通知を既読にする
func (h *NotificationHandler) MarkAllAsRead(c echo.Context) error {
	ctx := c.Request().Context()
	userID := context.GetCtxFromUser(ctx)

	updatedCount, err := h.notificationRepo.MarkAllAsRead(ctx, userID)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"updated_count": updatedCount,
		"message":       "All notifications marked as read",
	})
}

// DeleteNotification - 通知を削除
func (h *NotificationHandler) DeleteNotification(c echo.Context) error {
	ctx := c.Request().Context()
	userID := context.GetCtxFromUser(ctx)
	notificationID := c.Param("id")

	// 通知の所有権を確認
	notification, err := h.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Notification not found",
		})
	}

	if notification.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Access denied",
		})
	}

	if err := h.notificationRepo.Delete(ctx, notificationID); err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.NoContent(http.StatusNoContent)
}
