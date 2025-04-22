package dto

import (
	"mime/multipart"
)

type UpdateUserDTO struct {
	ID       uint64                `json:"id" binding:"required"`
	Name     string                `json:"name"`
	Profile  *multipart.FileHeader `form:"profile"`
	Password string                `form:"password"`
}

type GetUserDTO struct {
	ID uint64 `json:"id" binding:"required"`
}

type GetUserResponse struct {
	ID       uint64 `json:"id" binding:"required"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Profile  string `form:"profile"`
}
