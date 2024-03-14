package server

import (
	"fmt"

	"github.com/anchamber-studios/hevonen/lib/config"
	ldb "github.com/anchamber-studios/hevonen/lib/db"
	"github.com/anchamber-studios/hevonen/lib/events"
	"github.com/anchamber-studios/hevonen/lib/logger"
	m "github.com/anchamber-studios/hevonen/lib/middleware"
	"github.com/anchamber-studios/hevonen/services/admin/users/db"
	"github.com/anchamber-studios/hevonen/services/admin/users/services"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sqids/sqids-go"
)

type CustomContext struct {
	echo.Context
	Config       config.Config
	DB           *pgx.Conn
	IdConversion *sqids.Sqids
	Services     Services
}

type Services struct {
	Users *services.UserService
}

func Middleware(e *echo.Echo, conf config.Config) {
	// middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(m.Logging(logger.Get()))
	e.Use(customContext(conf))
}

func Routes(e *echo.Echo, conf config.Config) {
	unrestricted := e.Group("/users")
	unrestricted.POST("/login", login)
	unrestricted.POST("/register", new)

	restricted := e.Group("")
	restricted.Use(m.AuthPaseto(conf.TokenSecret))
	restricted.GET("/users", list)
	restricted.GET("/users/:userID", details)
	restricted.POST("/users/:userID/logout", logout)
}

func MiddlewareGroup(e *echo.Group, conf config.Config) {
	// middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(m.Logging(logger.Get()))
	e.Use(customContext(conf))
}

func RoutesGroup(e *echo.Group, conf config.Config) {
	unrestricted := e.Group("/users")
	unrestricted.POST("/login", login)
	unrestricted.POST("/register", new)

	restricted := e.Group("")
	restricted.Use(m.AuthPaseto(conf.TokenSecret))
	restricted.GET("/users", list)
	restricted.GET("/users/:userID", details)
	restricted.POST("/users/:userID/logout", logout)
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
			idc, err := setupIdConversion()
			if err != nil {
				c.Logger().Errorf("Unable setup id conversion: %v\n", err)
				c.Error(fmt.Errorf("internal server error"))
			}
			producer, err := events.NewEventProducerRedpanda(conf.Broker.Url)
			if err != nil {
				c.Logger().Errorf("Unable to connect to broker: %v\n", err)
				c.Error(fmt.Errorf("internal server error"))
			}
			repo := &db.UserRepoPostgre{
				DB:           conn,
				IdConversion: idc,
				Logger:       c.Logger(),
			}
			cc := &CustomContext{c, conf, conn, idc,
				Services{Users: services.NewUserService(repo, producer)},
			}

			return next(cc)
		}
	}
}

func setupIdConversion() (*sqids.Sqids, error) {
	s, err := sqids.New(sqids.Options{
		Alphabet: "FxnXM1kBN6cuhsAvjW3Co7l2RePyY8DwaU04Tzt9fHQrqSVKdpimLGIJOgb5ZE",
	})
	return s, err
}
