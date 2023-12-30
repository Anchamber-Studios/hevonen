package server

import (
	"fmt"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/services/users/db"
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
	Repos        Repos
}

type Repos struct {
	Users db.UserRepo
}

func Middleware(e *echo.Echo, conf config.Config) {
	// middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(customContext(conf))
}

func Routes(e *echo.Echo) {
	e.GET("/users", list)
	e.POST("/users/login", login)
	e.POST("/users/register", new)
	e.GET("/users/:userId", details)
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
			cc := &CustomContext{c, conf, conn, idc, Repos{
				Users: &db.UserRepoPostgre{
					DB:           conn,
					IdConversion: idc,
					Logger:       c.Logger(),
				},
			}}

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
