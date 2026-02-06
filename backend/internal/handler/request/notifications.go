package request

type MarkNotificationAsReadRequest struct {
	ID      string `param:"id" validate:"required,uuid"`
	Version int    `json:"version" validate:"required,min=1"`
}

type DeleteNotificationRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}
