package main

import (
	"context"
	"fmt"
	"log"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/services/general/profile/db"
	"github.com/anchamber-studios/hevonen/services/general/profile/server"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	log.Printf("Load .env file\n")
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	vars, err := godotenv.Read(".env")
	if err != nil {
		log.Println("Error reading .env file")
	}
	log.Printf("%v\n", vars)
	conf := config.LoadConfigWithVars(vars)

	e := echo.New()
	db.SetupDB(conf, e.Logger).Close(context.Background())
	server.Middleware(e, conf)
	server.Routes(e, conf)

	address := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	if conf.Tls.Enabled {
		e.Logger.Fatal(e.StartTLS(address, conf.Tls.Cert, conf.Tls.Key))
	} else {
		e.Logger.Fatal(e.Start(address))
	}
}
