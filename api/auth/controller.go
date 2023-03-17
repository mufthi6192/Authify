package auth

import "github.com/labstack/echo/v4"

func RegisterController(c echo.Context) error {
	response := Service(c).RegisterService()
	return c.JSON(response.Code, response)
}

func LoginController(c echo.Context) error {
	response := Service(c).LoginService()
	return c.JSON(response.Code, response)
}

func LogoutController(c echo.Context) error {
	response := Service(c).LogoutService()
	return c.JSON(response.Code, response)
}

func SendResetPasswordController(c echo.Context) error {
	response := Service(c).SendResetPasswordService()
	return c.JSON(response.Code, response)
}

func ResetPasswordController(c echo.Context) error {
	response := Service(c).ResetPasswordService()
	return c.JSON(response.Code, response)
}
