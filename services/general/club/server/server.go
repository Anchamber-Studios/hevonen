package server

import (
	"fmt"

	"github.com/anchamber-studios/hevonen/lib/config"
	repo "github.com/anchamber-studios/hevonen/services/club/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sqids/sqids-go"
)

type CustomContext struct {
	echo.Context
	Config       config.Config
	DB           *pgxpool.Conn
	IdConversion *sqids.Sqids
	Repos        Repos
}

type Repos struct {
	Members repo.MemberRepo
}

func Middleware(e *echo.Echo, conf config.Config) {

	// check db on startup
	pool := setupDb(conf, e.Logger)
	defer pool.Close()

	// middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(customContext(pool, conf))
}

func Routes(e *echo.Echo) {
	e.GET("/members", list)
	e.POST("/members", new)
	e.GET("/members/:memberId", details)
}

func customContext(pool *pgxpool.Pool, conf config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			conn, err := pool.Acquire(c.Request().Context())
			if err != nil {
				c.Logger().Errorf("Unable to connect to database: %v\n", err)
				c.Error(fmt.Errorf("internal server error"))
			}
			idc, err := setupIdConversion()
			if err != nil {
				c.Logger().Errorf("Unable setup id conversion: %v\n", err)
				c.Error(fmt.Errorf("internal server error"))
			}
			cc := &CustomContext{c, conf, conn, idc, Repos{
				Members: &repo.MemberRepoPostgre{DB: conn, IdConversion: idc},
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
