package middleware

import (
	responseFormatter "SMM-PPOB/helper/formatter"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func AdminAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user").(*jwt2.Token)
		claims := token.Claims.(jwt2.MapClaims)
		level := claims["level"].(string)

		if level != "admin" {
			response := responseFormatter.HttpResponse(401, "Gagal ! Anda tidak memiliki akses ke halaman admin", nil)
			return c.JSON(401, response)
		}

		return next(c)
	}
}
