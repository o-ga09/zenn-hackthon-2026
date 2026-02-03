package request

import "mime/multipart"

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

type CreateVLogRequest struct {
	Files       []*multipart.FileHeader `form:"files" validate:"required,min=1,dive,required"`
	MediaIDs    []string                `form:"mediaIds" validate:"omitempty,dive,uuid"`
	Title       *string                 `form:"title,omitempty"`
	TravelDate  *string                 `form:"travelDate,omitempty"`
	Destination *string                 `form:"destination,omitempty"`
	Theme       *string                 `form:"theme,omitempty"`
	MusicMood   *string                 `form:"musicMood,omitempty"`
	Duration    *int                    `form:"duration,omitempty"`
	Transition  *string                 `form:"transition,omitempty"`
}

type AnalyzeMediaRequest struct {
	Files    []*multipart.FileHeader `form:"files" validate:"required,min=1,dive,required"`
	MediaIDs []*string               `form:"mediaIds" validate:"omitempty,dive,uuid"`
}
