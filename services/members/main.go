package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sqids/sqids-go"
)

type TlsConfig struct {
	Enabled bool
	Key     string
	Cert    string
}
type Auth struct {
	ClientId     string
	ClientSecret string
}
type DB struct {
	Url      string
	Port     string
	Database string
	User     string
	Password string
}
type Config struct {
	Port string
	Host string
	Tls  TlsConfig
	Auth Auth
	DB   DB
}

func main() {
	// configuration
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	config := loadConfig()

	e := echo.New()

	// check db on startup
	pool := setupDb(config, e.Logger)
	defer pool.Close()

	// middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(customContext(pool, config))

	// handlers
	e.GET("/members", list)
	e.POST("/members", new)
	e.GET("/members/:memberId", details)

	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	if config.Tls.Enabled {
		e.Logger.Fatal(e.StartTLS(address, config.Tls.Cert, config.Tls.Key))
	} else {
		e.Logger.Fatal(e.Start(address))
	}
}

func loadConfig() Config {
	config := Config{
		Port: getOrDefault("PORT", "8443"),
		Host: getOrDefault("HOST", "[::1]"),
		Tls: TlsConfig{
			Enabled: false,
			Key:     getOrDefault("TLS_KEY", "certs/key.pem"),
			Cert:    getOrDefault("TLS_CERT", "certs/cert.pem"),
		},
		Auth: Auth{
			ClientId:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
		},
		DB: DB{
			Url:      os.Getenv("DB_URL"),
			Port:     os.Getenv("DB_PORT"),
			Database: os.Getenv("DB_DATABASE"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}

	if enabled, err := strconv.ParseBool(getOrDefault("TLS_ENABLED", "false")); err == nil && enabled {
		config.Tls.Enabled = true
	}
	return config
}

func setupIdConversion() (*sqids.Sqids, error) {
	s, err := sqids.New(sqids.Options{
		Alphabet: "FxnXM1kBN6cuhsAvjW3Co7l2RePyY8DwaU04Tzt9fHQrqSVKdpimLGIJOgb5ZE",
	})
	return s, err
}

type CustomContext struct {
	echo.Context
	Config       Config
	DB           *pgxpool.Conn
	IdConversion *sqids.Sqids
	Repos        Repos
}

type Repos struct {
	Members MemberRepo
}

func customContext(pool *pgxpool.Pool, config Config) echo.MiddlewareFunc {
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
			cc := &CustomContext{c, config, conn, idc, Repos{
				Members: &MemberRepoPostgre{DB: conn, IdConversion: idc},
			}}

			return next(cc)
		}
	}
}

func getOrDefault(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Printf("Environment variable %s not set, using default value %s\n", key, def)
		return def
	}
	return v
}
