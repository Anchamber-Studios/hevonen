package main

import (
	"fmt"
	"log"
	"os"

	"github.com/anchamber-studios/hevonen/frontend/pages/auth"
	"github.com/anchamber-studios/hevonen/frontend/pages/members"
	"github.com/anchamber-studios/hevonen/frontend/types"
	cclient "github.com/anchamber-studios/hevonen/services/club/client"
	uclient "github.com/anchamber-studios/hevonen/services/users/client"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := loadConfig()

	e := echo.New()

	// middleware
	e.Static("/public", "public")

	e.Use(middleware.CORS())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf",
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(session.Middleware((sessions.NewCookieStore([]byte(config.SessionSecret)))))
	e.Use(customContext(config))

	unrestricted := e.Group("/auth")
	unrestricted.GET("/login", auth.GetLogin)
	unrestricted.POST("/login", auth.PostLogin)
	unrestricted.GET("/register", auth.GetRegister)

	restricted := e.Group("")
	restricted.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(*types.CustomContext)
			if !cc.Session.LoggedIn {
				return c.Redirect(302, "/auth/login")
			}
			return next(c)
		}
	})
	restricted.GET("/", index)
	restricted.GET("/members", members.GetMemberList)
	restricted.GET("/members/new", members.GetNewMemberForm)
	restricted.POST("/members", members.GetNewMemberForm)

	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	if config.Tls.Enabled {
		log.Println("Starting server with TLS")
		e.Logger.Fatal(e.StartTLS(address, config.Tls.Cert, config.Tls.Key))
	} else {
		log.Println("Starting server without TLS")
		e.Logger.Fatal(e.Start(address))
	}
}

func loadConfig() types.Config {
	// configuration
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	return types.Config{
		Host:          getOrDefault(os.Getenv("HOST"), "[::0]"),
		Port:          getOrDefault(os.Getenv("PORT"), "4443"),
		Clients:       createClients(),
		SessionSecret: getOrDefault(os.Getenv("SESSION_SECRET"), "session_secret"),
	}
}

func customContext(config types.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &types.CustomContext{
				Context: c,
				Config:  config,
				Session: types.Session{LoggedIn: false},
			}
			sess, _ := session.Get("session", c)
			if sess != nil {
				if val, ok := sess.Values["id"].(string); ok {
					cc.Session.ID = val
				}
				if val, ok := sess.Values["email"].(string); ok {
					cc.Session.Email = val
				}
				cc.Session.LoggedIn = cc.Session.ID != "" && cc.Session.Email != ""
			}
			return next(cc)
		}
	}
}

func createClients() types.Clients {
	return types.Clients{
		Members: &cclient.MemberClient{
			Url: getOrDefault(os.Getenv("MEMBERS_URL"), "http://localhost:8443/members"),
		},
		User: &uclient.UserClient{
			Url: getOrDefault(os.Getenv("USERS_URL"), "http://localhost:7443/users"),
		},
	}
}

func getOrDefault(s string, d string) string {
	if s == "" {
		return d
	}
	return s
}
