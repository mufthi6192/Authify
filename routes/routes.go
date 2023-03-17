package routes

import (
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) error {

	g := e.Group("/api") //Use prefix for API Routes

	//Define API Routes
	err := ApiRoutes(g).Api()

	if err != nil {
		panic(err)
	}

	//Define Web Routes
	err = WebRoutes(g).Web()

	if err != nil {
		panic(err)
	}

	return nil
}
