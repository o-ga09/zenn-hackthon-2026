package handler

import "github.com/labstack/echo"

type IUserServer interface {
	GET(ctx echo.Context) error
	Create(ctx echo.Context) error
	Update(ctx echo.Context) error
	Delete(ctx echo.Context) error
}

type UserServer struct{}

func (u *UserServer) GET(ctx echo.Context) error {
	return ctx.String(200, "Get User")
}

func (u *UserServer) Create(ctx echo.Context) error {
	return ctx.String(200, "Create User")
}

func (u *UserServer) Update(ctx echo.Context) error {
	return ctx.String(200, "Update User")
}

func (u *UserServer) Delete(ctx echo.Context) error {
	return ctx.String(200, "Delete User")
}
