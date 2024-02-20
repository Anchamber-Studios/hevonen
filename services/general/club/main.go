package main

import (
	"fmt"
	"log"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/services/club/db"
	"github.com/anchamber-studios/hevonen/services/club/server"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	conf := config.LoadConfig()

	e := echo.New()
	if err = db.SetupDB(conf); err != nil {
		log.Fatalf("Unable to setup database: %v\n", err)
	}
	server.Middleware(e, conf)
	server.Routes(e)

	address := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	if conf.Tls.Enabled {
		e.Logger.Fatal(e.StartTLS(address, conf.Tls.Cert, conf.Tls.Key))
	} else {
		e.Logger.Fatal(e.Start(address))
	}
}
