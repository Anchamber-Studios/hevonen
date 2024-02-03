package config

import (
	"log"
	"os"
	"strconv"

	"github.com/lestrrat-go/jwx/v2/jwk"
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
	TokenSecret string
	JWK         jwk.Set
	Tls         TlsConfig
	Auth        Auth
	DB          DB
	Broker      Broker
}

func LoadConfig() Config {
	config := Config{
		Port:        getOrDefault("PORT", "8443"),
		Host:        getOrDefault("HOST", "localhost"),
		TokenSecret: getOrDefault("TOKEN_SECRET", "1234567890123456789012345678901212345678901234567890123456789012"),
		JWK:         NewJWKSet(os.Getenv("JWK_FILE")),
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

func LoadConfigWithVars(vars map[string]string) Config {
	config := Config{
		Port:        getOrDefaultWithVars(vars, "PORT", "8443"),
		Host:        getOrDefaultWithVars(vars, "HOST", "localhost"),
		TokenSecret: getOrDefaultWithVars(vars, "TOKEN_SECRET", "1234567890123456789012345678901212345678901234567890123456789012"),
		JWK:         NewJWKSet(getOsOrVars(vars, "JWK_FILE")),
		Tls: TlsConfig{
			Enabled: false,
			Key:     getOrDefaultWithVars(vars, "TLS_KEY", "certs/key.pem"),
			Cert:    getOrDefaultWithVars(vars, "TLS_CERT", "certs/cert.pem"),
		},
		Auth: Auth{
			ClientId:     getOsOrVars(vars, "CLIENT_ID"),
			ClientSecret: getOsOrVars(vars, "CLIENT_SECRET"),
		},
		DB: DB{
			Url:      getOsOrVars(vars, "DB_URL"),
			Port:     getOsOrVars(vars, "DB_PORT"),
			Database: getOsOrVars(vars, "DB_DATABASE"),
			User:     getOsOrVars(vars, "DB_USER"),
			Password: getOsOrVars(vars, "DB_PASSWORD"),
		},
		Broker: Broker{
			Url: getOrDefaultWithVars(vars, "BROKER_URL", "localhost:19092"),
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

func getOrDefaultWithVars(vars map[string]string, key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		if val, ok := vars[key]; ok {
			return val
		}

		log.Printf("Environment variable %s not set, using default value %s\n", key, def)
		return def
	}
	return v
}

func getOsOrVars(vars map[string]string, key string) string {
	v := os.Getenv(key)
	if v == "" {
		if val, ok := vars[key]; ok {
			return val
		}

		log.Panicf("Variable %s is not set. Cannot start server.\n", key)
	}
	return v
}

func NewJWKSet(path string) jwk.Set {
	f, err := os.Open(path)
	if err != nil {
		panic("failed to open jwk file")
	}
	defer f.Close()
	jwkSet, err := jwk.ParseReader(f)
	if err != nil {
		panic("failed to parse jwk file")
	}
	public, err := jwk.PublicSetOf(jwkSet)
	if err != nil {
		panic("failed to parse public keys of jwk file")
	}
	return public
}
