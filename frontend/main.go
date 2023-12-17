package main

import (
	"fmt"
	"log"
	"os"

	"github.com/anchamber-studios/hevonen/services/members/client"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Host    string
	Port    string
	Tls     TlsConfig
	Clients Clients
}

type TlsConfig struct {
	Enabled bool
	Key     string
	Cert    string
}

type Clients struct {
	Members *client.MemberClient
}

func main() {
	config := loadConfig()

	e := echo.New()

	// middleware
	e.Static("/public", "public")

	e.Use(middleware.CORS())
	e.Use(middleware.CSRF())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(customContext(config))

	e.GET("/", index)
	e.GET("/members", memberList)
	e.POST("/members", postNewMember)
	e.GET("/members/new", memberNew)

	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	if config.Tls.Enabled {
		e.Logger.Fatal(e.StartTLS(address, config.Tls.Cert, config.Tls.Key))
	} else {
		e.Logger.Fatal(e.Start(address))
	}
}

func loadConfig() Config {
	// configuration
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	return Config{
		Host:    getOrDefault(os.Getenv("HOST"), "localhost"),
		Port:    getOrDefault(os.Getenv("PORT"), "4443"),
		Clients: createClients(),
	}
}

type CustomContext struct {
	echo.Context
	Config Config
}

func customContext(config Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			cc := &CustomContext{c, config}

			return next(cc)
		}
	}
}

func createClients() Clients {
	return Clients{
		Members: &client.MemberClient{
			Url: getOrDefault(os.Getenv("MEMBERS_URL"), "http://localhost:8443/members"),
		},
	}
}

func getOrDefault(s string, d string) string {
	if s == "" {
		return d
	}
	return s
}
