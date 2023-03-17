package routes

import "github.com/labstack/echo/v4"

type routeWeb struct {
	G *echo.Group
}

func WebRoutes(g *echo.Group) WebRouteInterface {
	return &routeWeb{
		G: g,
	}
}

func (r *routeWeb) Web() error {

	//Initialize *echo.Group
	//todo:Implement *echo.Group

	//Define Global Middleware here
	//todo:Implement Middleware

	//Register web routes here
	//todo:Register Routes

	return nil
}
