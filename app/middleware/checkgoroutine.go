package middleware

import (
	"github.com/labstack/echo/v4"
	"log"
	"runtime"
)

func RunningGoroutine(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("Running Goroutine :", runtime.NumGoroutine())
		return next(c)
	}
}
