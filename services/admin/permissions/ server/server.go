package server

import (
	"github.com/anchamber-studios/hevonen/lib/config"
	m "github.com/anchamber-studios/hevonen/lib/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Middleware(e *echo.Echo, conf config.Config) {
	// middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(m.Logging())
	e.Use(customContext(conf))
}

func Routes(e *echo.Echo, conf config.Config) {
	restricted := e.Group("")
	restricted.Use(m.AuthPaseto(conf.TokenSecret))
}

type CustomContext struct {
	echo.Context
	Config config.Config
}

func customContext(conf config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{
				Context: c,
				Config:  conf,
			}
			return next(cc)
		}
	}
}
