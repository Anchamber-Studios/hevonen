package middleware

import (
	"github.com/anchamber-studios/hevonen/lib/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func AuthPaseto(tokenKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "missing authorization header")
			}
			// _ := auth[len("Bearer "):]
			return next(c)
		}
	}
}

func AuthOry() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil || cookie == nil {
				c.Logger().Infof("missing session cookie")
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "missing session cookie")
			}

			// validate the session cookie with ory

			// Continue with the next middleware or handler
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
				zap.String("method", v.Method),
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	})
}
