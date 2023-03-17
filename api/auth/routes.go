package auth

import (
	"SMM-PPOB/app/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Routes(g *echo.Group) {

	//Initialize *echo.Group
	r := g.Group("/auth")

	//Define Middleware Here

	//Define Routes Here
	r.POST("/send-reset-password", SendResetPasswordController)
	r.POST("/reset-password", ResetPasswordController)
	r.POST("/register", RegisterController)
	r.POST("/login", LoginController)
	r.Use(echojwt.JWT([]byte("ice_dolce_latte")))
	r.Use(middleware.Authentication)
	r.GET("/logout", LogoutController)
}
