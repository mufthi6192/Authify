package profile

import (
	"SMM-PPOB/app/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Routes(g *echo.Group) {

	//Initialize *echo.Group
	r := g.Group("/client/profile")
	nr := g.Group("/client/verification")

	//Define Routes without using Middleware Here
	nr.GET("/email", UpdateEmailVerificationStatus)

	//Define Routes using Middleware Here
	nr.Use(echojwt.JWT([]byte("ice_dolce_latte")))
	nr.Use(middleware.Authentication)
	r.Use(echojwt.JWT([]byte("ice_dolce_latte")))
	r.Use(middleware.Authentication)

	r.GET("/", GetProfileController)

	r.GET("/login-history", GetLoginHistoryController)
	r.GET("/latest-login-history", GetLatestLoginHistory)

	r.PATCH("/change-password", ChangePasswordController)

	nr.GET("/resend-email", ResendVerificationEmailController)

}
