package main

import (
	"fmt"
	"log"

	"github.com/anchamber-studios/hevonen/services/members/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sqids/sqids-go"
)

func main() {
	// configuration
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	conf := config.LoadConfig()

	e := echo.New()

	// check db on startup
	pool := setupDb(conf, e.Logger)
	defer pool.Close()

	// middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(customContext(pool, conf))

	// handlers
	e.GET("/members", list)
	e.POST("/members", new)
	e.GET("/members/:memberId", details)

	address := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	if conf.Tls.Enabled {
		e.Logger.Fatal(e.StartTLS(address, conf.Tls.Cert, conf.Tls.Key))
	} else {
		e.Logger.Fatal(e.Start(address))
	}
}

func setupIdConversion() (*sqids.Sqids, error) {
	s, err := sqids.New(sqids.Options{
		Alphabet: "FxnXM1kBN6cuhsAvjW3Co7l2RePyY8DwaU04Tzt9fHQrqSVKdpimLGIJOgb5ZE",
	})
	return s, err
}

type CustomContext struct {
	echo.Context
	Config       config.Config
	DB           *pgxpool.Conn
	IdConversion *sqids.Sqids
	Repos        Repos
}

type Repos struct {
	Members MemberRepo
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
				Members: &MemberRepoPostgre{DB: conn, IdConversion: idc},
			}}

			return next(cc)
		}
	}
}
