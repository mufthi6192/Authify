package profile

import (
	"SMM-PPOB/app/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Routes(g *echo.Group) {

	//Initialize *echo.Group
	r := g.Group("/admin/profile")

	//Define Middleware Here

	//Define Routes Here
	r.Use(echojwt.JWT([]byte("ice_dolce_latte")))
	r.Use(middleware.Authentication)
	r.Use(middleware.AdminAccess)

	r.GET("/", GetProfileController)

	r.GET("/login-history", GetLoginHistoryController)
	r.GET("/latest-login-history", GetLatestLoginHistory)

	r.PATCH("/change-password", ChangePasswordController)
}
