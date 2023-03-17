package profile

import "github.com/labstack/echo/v4"

func GetProfileController(c echo.Context) error {
	response := Service(c).GetProfileService()
	return c.JSON(response.Code, response)
}

func GetLoginHistoryController(c echo.Context) error {
	response := Service(c).GetLoginHistoryService()
	return c.JSON(response.Code, response)
}

func GetLatestLoginHistory(c echo.Context) error {
	response := Service(c).GetLatestLoginHistoryService()
	return c.JSON(response.Code, response)
}

func ChangePasswordController(c echo.Context) error {
	response := Service(c).ChangePasswordService()
	return c.JSON(response.Code, response)
}
