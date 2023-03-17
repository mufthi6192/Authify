package test

import "github.com/labstack/echo/v4"

func Routes(g *echo.Group) {

	//Initialize *echo.Group
	r := g.Group("/test")

	//Define Middleware Here

	//Define Routes Here
	r.GET("", OkController)
}
