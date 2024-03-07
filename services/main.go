package main

import (
	"fmt"
	"log"

	"github.com/anchamber-studios/hevonen/lib/config"

	authdb "github.com/anchamber-studios/hevonen/services/admin/auth/db"
	authserver "github.com/anchamber-studios/hevonen/services/admin/auth/server"
	clubdb "github.com/anchamber-studios/hevonen/services/club/db"
	clubserver "github.com/anchamber-studios/hevonen/services/club/server"
	profiledb "github.com/anchamber-studios/hevonen/services/general/profile/db"
	profileserver "github.com/anchamber-studios/hevonen/services/general/profile/server"
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
	conf := config.LoadConfigWithVars(vars)

	e := echo.New()
	setupDomains(e, conf)

	address := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	if conf.Tls.Enabled {
		e.Logger.Fatal(e.StartTLS(address, conf.Tls.Cert, conf.Tls.Key))
	} else {
		e.Logger.Fatal(e.Start(address))
	}
}

func setupDomains(e *echo.Echo, conf config.Config) {
	go setupDB(conf)
	log.Printf("AuthService: Start setup\n")
	usersPath := fmt.Sprintf("/%s", "/api/auth")
	usersGroup := e.Group(usersPath)
	authserver.MiddlewareGroup(usersGroup, conf)
	authserver.RoutesGroup(usersGroup, conf)
	log.Printf("AuthService: Finished setup\n")

	log.Printf("Start setup of club service\n")
	clubPath := fmt.Sprintf("/%s", "/api/club")
	clubGroup := e.Group(clubPath)
	clubserver.MiddlewareGroup(clubGroup, conf)
	clubserver.RoutesGroup(clubGroup)
	log.Printf("Finished club of profile service\n")

	log.Printf("Start setup of profile service\n")
	profilePath := fmt.Sprintf("/%s", "/api/profile")
	profileGroup := e.Group(profilePath)
	profileserver.MiddlewareGroup(profileGroup, conf)
	profileserver.RoutesGroup(profileGroup)
	log.Printf("Finished setup of profile service\n")

}

func setupDB(conf config.Config) {
	log.Printf("AuthService: setup db\n")
	if err := authdb.SetupDB(conf); err != nil {
		log.Fatalf("AuthService: Unable to setup database: %v\n", err)
	}
	log.Printf("AuthService: finished setup db\n")

	log.Printf("ProfileService: setup db\n")
	if err := clubdb.SetupDB(conf); err != nil {
		log.Fatalf("ProfileService: Unable to setup database: %v\n", err)
	}
	log.Printf("ProfileService: finished setup db\n")

	log.Printf("ProfileService: setup db\n")
	if err := profiledb.SetupDB(conf); err != nil {
		log.Fatalf("ProfileService: Unable to setup database: %v\n", err)
	}
	log.Printf("ProfileService: finished setup db\n")

}
