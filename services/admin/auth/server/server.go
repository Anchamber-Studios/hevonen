package server

import (
	"context"
	"fmt"

	"github.com/anchamber-studios/hevonen/lib/config"
	ldb "github.com/anchamber-studios/hevonen/lib/db"
	"github.com/anchamber-studios/hevonen/lib/events"
	"github.com/anchamber-studios/hevonen/lib/logger"
	m "github.com/anchamber-studios/hevonen/lib/middleware"
	"github.com/anchamber-studios/hevonen/services/admin/auth/db"
	"github.com/anchamber-studios/hevonen/services/admin/auth/services"
	us "github.com/anchamber-studios/hevonen/services/admin/users/services"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Middleware(e *echo.Echo, conf config.Config) {
	// middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(m.Logging(logger.Get()))
	e.Use(customContext(conf))
}

func Routes(e *echo.Echo, conf config.Config) {
	restricted := e.Group("")
	restricted.Use(m.AuthPaseto(conf.TokenSecret))
	restricted.GET("/groups", GetGroups)

	restricted.GET("/auth", GetAuthorizations)

	restricted.GET("/services", GetServices)
	restricted.GET("/services/:serviceId/auth", GetAuthorizationsForService)
}

func Events(e *echo.Echo, conf config.Config) {
	topic := us.GetTopicName(us.ActionLogin)
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
	DB       *pgx.Conn
	Config   config.Config
	Services Services
}

type Services struct {
	AuthService services.AuthService
}

func customContext(conf config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			conn, err := ldb.OpenConnection(conf)
			if err != nil {
				c.Logger().Errorf("Unable to connect to database: %v\n", err)
				c.Error(fmt.Errorf("internal server error"))
			}
			defer conn.Close(c.Request().Context())

			cc := &CustomContext{
				DB:      conn,
				Context: c,
				Config:  conf,
				Services: Services{
					AuthService: &services.AuthServiceImpl{
						AuthRepo: &db.AuthRepositoryPostgre{},
					},
				},
			}
			return next(cc)
		}
	}
}
