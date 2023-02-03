package middleware

import (
	"campyuk-api/config"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JWTWithConfig() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte(config.JWT_KEY),
		SigningMethod: jwt.SigningMethodHS256.Name,
		ErrorHandlerWithContext: func(err error, c echo.Context) error {
			// return c.JSON(401, map[string]interface{}{"message": "access is denied due to invalid credential"})
			return nil
		},
		ContinueOnIgnoredError: true,
	})
}
