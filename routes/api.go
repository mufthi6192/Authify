package routes

import (
	"SMM-PPOB/api/admin/profile"
	"SMM-PPOB/api/admin/user"
	"SMM-PPOB/api/auth"
	clientProfile "SMM-PPOB/api/client/profile"
	"SMM-PPOB/api/test"
	middleware2 "SMM-PPOB/app/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type routeApi struct {
	G *echo.Group
}

func ApiRoutes(e *echo.Group) ApiRouteInterface {
	return &routeApi{
		G: e,
	}
}

func (r *routeApi) Api() error {

	//Initialize *echo.Group
	g := r.G

	//Define Global Middleware here
	g.Use(middleware.CORS())
	g.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))
	g.Use(middleware2.RunningGoroutine)

	//Register api routes here
	test.Routes(g)
	auth.Routes(g)
	profile.Routes(g)
	clientProfile.Routes(g)
	user.Routes(g)

	return nil
}
