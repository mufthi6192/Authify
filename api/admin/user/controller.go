package user

import "github.com/labstack/echo/v4"

func AddUserController(ctx echo.Context) error {
	response := Service(ctx).AddUserService()
	return ctx.JSON(response.Code, response)
}

func DeleteUserController(ctx echo.Context) error {
	response := Service(ctx).DeleteUserService()
	return ctx.JSON(response.Code, response)
}
