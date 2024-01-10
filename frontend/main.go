package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anchamber-studios/hevonen/frontend/pages/admin"
	"github.com/anchamber-studios/hevonen/frontend/pages/auth"
	"github.com/anchamber-studios/hevonen/frontend/pages/general/profile"
	"github.com/anchamber-studios/hevonen/frontend/pages/members"
	"github.com/anchamber-studios/hevonen/frontend/types"
	m "github.com/anchamber-studios/hevonen/lib/middleware"
	uclient "github.com/anchamber-studios/hevonen/services/admin/users/client"
	cclient "github.com/anchamber-studios/hevonen/services/club/client"
	pclient "github.com/anchamber-studios/hevonen/services/general/profile/client"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	ory "github.com/ory/client-go"
)

func main() {
	config := loadConfig()

	e := echo.New()

	// middleware
	e.Static("/public", "public")

	e.Use(middleware.CORS())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "cookie:_csrf",
		CookiePath:     "/",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())
	e.Use(m.Logging())
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
	restricted.GET("/auth/logout", auth.GetLogout)

	restricted.GET("/", index)
	restricted.GET("/members", members.GetMemberList)
	restricted.GET("/members/new", members.GetNewMemberForm)
	restricted.POST("/members", members.GetNewMemberForm)

	restricted.GET("/admin/users", admin.GetUsers)
	restricted.GET("/admin/users/:userId", admin.GetUser)

	restricted.GET("/u/p", profile.GetProfile)
	restricted.PUT("/u/p", profile.UpdateProfile)

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
		Host:          getOrDefault("HOST", "[::0]"),
		Port:          getOrDefault("PORT", "4443"),
		Clients:       createClients(),
		SessionSecret: getOrDefault("SESSION_SECRET", "session_secret"),
		TokenSecret:   getOrDefault("TOKEN_SECRET", "1234567890123456789012345678901212345678901234567890123456789012"),
	}
}

func customContext(config types.Config) echo.MiddlewareFunc {
	c := ory.NewConfiguration()
	c.Servers = ory.ServerConfigurations{{URL: fmt.Sprintf("http://localhost:%s/.ory", "4000")}}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cc := &types.CustomContext{
				Context:   ctx,
				Auth:      ory.NewAPIClient(c),
				Config:    config,
				Session:   types.Session{LoggedIn: false},
				HXRequest: ctx.Request().Header.Get(HX_REQUEST_HEADER) == "true",
			}

			// Get the session cookie from the request
			cookie, err := ctx.Cookie("session")
			if err != nil || cookie == nil {
				return next(cc)
			}
			// todo how to tokenize session token
			// https://www.ory.sh/docs/identities/session-to-jwt-cors
			session, _, err := cc.Auth.FrontendAPI.ToSession(ctx.Request().Context()).
				XSessionToken(cookie.Value).
				TokenizeAs("jwt_example_template").
				Execute()
			if err != nil {
				cc.Logger().Errorf("Unable to get session: %v\n", err)
				return next(cc)
			}
			cc.Logger().Infof("session: %v\n", session)
			if session != nil {
				cc.Session.Token = *session.Tokenized
				if traits, ok := session.Identity.Traits.(map[string]interface{}); ok {
					cc.Session.Email = traits["email"].(string)
				}
				cc.Session.ID = session.Identity.Id
				cc.Session.LoggedIn = true
			}

			return next(cc)
		}
	}
}

func createClients() types.Clients {
	return types.Clients{
		Members: cclient.MemberClient{
			Url: getOrDefault("MEMBERS_URL", "http://localhost:8443/members"),
		},
		User: &uclient.UserClientHttp{
			Url: getOrDefault("USERS_URL", "http://localhost:7000/users"),
		},
		Profile: &pclient.ProfileClientHttp{
			Url: getOrDefault("PROFILE_URL", "http://localhost:7002/profiles"),
		},
	}
}

func getOrDefault(s string, d string) string {
	v := os.Getenv(s)
	if v == "" {
		return d
	}
	return s
}
