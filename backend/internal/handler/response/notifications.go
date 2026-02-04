package response

import (
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/date"
)

type Notification struct {
	ID        string `json:"id"`
	Version   int    `json:"version"`
	UserID    string `json:"user_id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	MediaID   string `json:"media_id,omitempty"`
	VlogID    string `json:"vlog_id,omitempty"`
	Read      bool   `json:"read"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type NotificationListResponse struct {
	Notifications []*Notification `json:"notifications"`
	UnreadCount   int             `json:"unread_count"`
}

func ToNotificationResponse(n []*domain.Notification, unreadCount int) *NotificationListResponse {
	res := make([]*Notification, 0, len(n))
	for _, notif := range n {
		res = append(res, ToNotification(notif))
	}
	return &NotificationListResponse{
		Notifications: res,
		UnreadCount:   unreadCount,
	}
}

func ToNotification(n *domain.Notification) *Notification {
	return &Notification{
		ID:        n.ID,
		Version:   n.Version,
		UserID:    n.UserID,
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		MediaID:   n.MediaID.String,
		VlogID:    n.VlogID.String,
		Read:      n.Read,
		CreatedAt: date.Format(n.CreatedAt),
		UpdatedAt: date.Format(n.UpdatedAt),
	}
}
