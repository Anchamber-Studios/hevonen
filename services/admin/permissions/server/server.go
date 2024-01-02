package server

import (
	"context"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/lib/events"
	"github.com/anchamber-studios/hevonen/lib/logger"
	m "github.com/anchamber-studios/hevonen/lib/middleware"
	"github.com/anchamber-studios/hevonen/services/admin/users/services"
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

func Events(e *echo.Echo, conf config.Config) {
	topic := services.GetTopicName(services.ActionLogin)
	loginEvents, err := events.NewEventConsumerRedpanda([]string{conf.Broker.Url}, topic)
	if err != nil {
		e.Logger.Fatal(err)
	}
	log := logger.Get()
	ctx := logger.WithCtx(context.Background(), log)
	log.Sugar().Infof("Subscribing to topic '%s'\n", topic)
	loginEvents.Subscribe(ctx, func(eventCtx context.Context, data []byte, headers map[string]string) error {
		log := logger.FromContext(eventCtx)
		log.Info("Received login event")
		return nil
	})
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
