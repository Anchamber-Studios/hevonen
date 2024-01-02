package config

import (
	"log"
	"os"
	"strconv"
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
type Broker struct {
	Url string
}
type Config struct {
	Port        string
	Host        string
	Tls         TlsConfig
	Auth        Auth
	DB          DB
	TokenSecret string
	Broker      Broker
}

func LoadConfig() Config {
	config := Config{
		Port:        getOrDefault("PORT", "8443"),
		Host:        getOrDefault("HOST", "[::1]"),
		TokenSecret: getOrDefault("TOKEN_SECRET", "1234567890123456789012345678901212345678901234567890123456789012"),
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
		Broker: Broker{
			Url: getOrDefault("BROKER_URL", "localhost:19092"),
		},
	}

	if enabled, err := strconv.ParseBool(getOrDefault("TLS_ENABLED", "false")); err == nil && enabled {
		config.Tls.Enabled = true
	}
	return config
}

func getOrDefault(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Printf("Environment variable %s not set, using default value %s\n", key, def)
		return def
	}
	return v
}
