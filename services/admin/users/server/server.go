package server

import (
	"fmt"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/lib/events"
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
	e.Use(m.Logging())
	e.Use(customContext(conf))
}

func Routes(e *echo.Echo, conf config.Config) {
	unrestricted := e.Group("/users")
	unrestricted.POST("/login", login)
	unrestricted.POST("/register", new)

	restricted := e.Group("")
	restricted.Use(m.AuthPaseto(conf.TokenSecret))
	restricted.GET("/users", list)
	restricted.GET("/users/:userId", details)
	restricted.POST("/users/:userId/logout", logout)
}

func customContext(conf config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			conn, err := db.OpenConnection(conf, c.Logger())
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
