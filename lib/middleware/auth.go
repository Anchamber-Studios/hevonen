package middleware

import (
	"aidanwoods.dev/go-paseto"
	"github.com/labstack/echo/v4"
)

func AuthPaseto(tokenKey string) echo.MiddlewareFunc {
	key, err := paseto.V4SymmetricKeyFromHex(tokenKey)
	parser := paseto.NewParser()
	if err != nil {
		panic(err)
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "missing authorization header")
			}
			signed := auth[len("Bearer "):]
			_, err := parser.ParseV4Local(key, signed, nil)
			if err != nil {
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "invalid token")
			}
			return next(c)
		}
	}
}
