package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/request"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/response"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
)

type INotificationHandler interface {
	GetNotifications(ctx echo.Context) error
	MarkAsRead(ctx echo.Context) error
	MarkAllAsRead(ctx echo.Context) error
	DeleteNotification(ctx echo.Context) error
	DeleteAllNotifications(ctx echo.Context) error
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

	return c.JSON(http.StatusOK, response.ToNotificationResponse(notifications, int(unreadCount)))
}

// MarkAsRead - 通知を既読にする
func (h *NotificationHandler) MarkAsRead(c echo.Context) error {
	ctx := c.Request().Context()
	userID := context.GetCtxFromUser(ctx)

	var req request.MarkNotificationAsReadRequest
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}
	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	// 通知の所有権を確認
	notification, err := h.notificationRepo.FindByID(ctx, req.ID)
	if err != nil {
		return errors.MakeNotFoundError(ctx, "通知が取得できませんでした")
	}

	if notification.UserID != userID {
		return errors.MakeBusinessError(ctx, "ユーザーIDが不正です")
	}

	if err := h.notificationRepo.MarkAsRead(ctx, &domain.Notification{BaseModel: domain.BaseModel{ID: req.ID, Version: req.Version}}); err != nil {
		return errors.Wrap(ctx, err)
	}

	// 更新後の通知を返す
	updatedNotification, err := h.notificationRepo.FindByID(ctx, req.ID)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.JSON(http.StatusOK, response.ToNotification(updatedNotification))
}

// MarkAllAsRead - 全通知を既読にする
func (h *NotificationHandler) MarkAllAsRead(c echo.Context) error {
	ctx := c.Request().Context()
	userID := context.GetCtxFromUser(ctx)

	_, err := h.notificationRepo.MarkAllAsRead(ctx, &domain.Notification{UserID: userID})
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteNotification - 通知を削除
func (h *NotificationHandler) DeleteNotification(c echo.Context) error {
	ctx := c.Request().Context()
	userID := context.GetCtxFromUser(ctx)
	var req request.DeleteNotificationRequest
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}
	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	// 通知の所有権を確認
	notification, err := h.notificationRepo.FindByID(ctx, req.ID)
	if err != nil {
		return errors.MakeNotFoundError(ctx, "通知が取得できませんでした")
	}

	if notification.UserID != userID {
		return errors.MakeBusinessError(ctx, "ユーザーIDが不正です")
	}

	if err := h.notificationRepo.Delete(ctx, notification); err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteAllNotifications - 全通知を削除
func (h *NotificationHandler) DeleteAllNotifications(c echo.Context) error {
	ctx := c.Request().Context()
	userID := context.GetCtxFromUser(ctx)

	if err := h.notificationRepo.DeleteAllByUserID(ctx, userID); err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.NoContent(http.StatusNoContent)
}
