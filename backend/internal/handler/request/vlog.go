package request

type VLogListRequest struct {
	Offset *int `query:"offset" validate:"omitempty"`
	Limit  *int `query:"limit" validate:"omitempty"`
}

type VLogGetByIDRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

type VLogDeleteRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}
