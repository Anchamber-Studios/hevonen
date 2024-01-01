package middleware

import (
	"aidanwoods.dev/go-paseto"
	"github.com/anchamber-studios/hevonen/lib/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
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

func Logging() echo.MiddlewareFunc {
	logger := logger.Get()
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	})
}
