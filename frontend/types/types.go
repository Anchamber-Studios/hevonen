package types

import (
	cclient "github.com/anchamber-studios/hevonen/services/club/client"
	uclient "github.com/anchamber-studios/hevonen/services/users/client"
	"github.com/labstack/echo/v4"
)

type Session struct {
	LoggedIn bool
	ID       string
	Email    string
}

type CustomContext struct {
	echo.Context
	Config  Config
	Session Session
}

type Config struct {
	Host          string
	Port          string
	Tls           TlsConfig
	Clients       Clients
	SessionSecret string
}

type TlsConfig struct {
	Enabled bool
	Key     string
	Cert    string
}

type Clients struct {
	Members *cclient.MemberClient
	User    *uclient.UserClient
}
