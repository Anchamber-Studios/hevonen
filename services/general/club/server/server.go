package server

import (
	"fmt"

	"github.com/anchamber-studios/hevonen/lib/config"
	ldb "github.com/anchamber-studios/hevonen/lib/db"
	"github.com/anchamber-studios/hevonen/lib/logger"
	m "github.com/anchamber-studios/hevonen/lib/middleware"
	"github.com/anchamber-studios/hevonen/services/club/db"
	"github.com/anchamber-studios/hevonen/services/club/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sqids/sqids-go"
)

type CustomContext struct {
	echo.Context
	Config   config.Config
	Services Services
}

type Services struct {
	Clubs   *services.ClubService
	Members *services.MemberService
}

func Middleware(e *echo.Echo, conf config.Config) {

	// middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(m.Logging(logger.Get()))
	e.Use(customContext(conf))
	e.Use(m.AuthJWTOry(conf.JWK))
}

func Routes(e *echo.Echo) {
	clubHandler := &ClubHandler{}
	e.GET("/c", clubHandler.ListForIdentity).Name = "ListForIdentity"
	e.POST("/c", clubHandler.Create).Name = "CreateClub"
}

func MiddlewareGroup(e *echo.Group, conf config.Config) {

	// middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(m.Logging(logger.Get()))
	e.Use(customContext(conf))
	e.Use(m.AuthJWTOry(conf.JWK))
}

func RoutesGroup(e *echo.Group) {
	clubHandler := &ClubHandler{}
	e.GET("/c", clubHandler.ListForIdentity).Name = "ListForIdentity"
	e.POST("/c", clubHandler.Create).Name = "CreateClub"
}

func customContext(conf config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			conn, err := ldb.OpenConnection(conf)
			if err != nil {
				c.Logger().Errorf("Unable to connect to database: %v\n", err)
				c.Error(fmt.Errorf("internal server error"))
			}
			idc, err := SetupIdConversion()
			if err != nil {
				c.Logger().Errorf("Unable setup id conversion: %v\n", err)
				c.Error(fmt.Errorf("internal server error"))
			}
			cc := &CustomContext{c, conf,
				Services{
					Clubs:   services.NewClubService(&db.ClubRepoPostgre{DB: conn, IdConversion: idc}),
					Members: services.NewMemberService(&db.MemberRepoPostgre{DB: conn, IdConversion: idc}),
				},
			}

			return next(cc)
		}
	}
}

func SetupIdConversion() (*sqids.Sqids, error) {
	s, err := sqids.New(sqids.Options{
		Alphabet: "FxnXM1kBN6cuhsAvjW3Co7l2RePyY8DwaU04Tzt9fHQrqSVKdpimLGIJOgb5ZE",
	})
	return s, err
}
