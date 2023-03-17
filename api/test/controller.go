package test

import "github.com/labstack/echo/v4"

func OkController(c echo.Context) error {
	return c.JSON(200, "OK")
}
