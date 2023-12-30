package types

import (
	"github.com/anchamber-studios/hevonen/lib"
	cclient "github.com/anchamber-studios/hevonen/services/club/client"
	uclient "github.com/anchamber-studios/hevonen/services/users/client"
	"github.com/labstack/echo/v4"
)

type Session struct {
	LoggedIn bool
	ID       string
	Email    string
	Token    string
}

type CustomContext struct {
	echo.Context
	Config    Config
	Session   Session
	HXRequest bool
}

func (cc *CustomContext) ClientContext() lib.ClientContext {
	return lib.NewClientContext(cc.Context, cc.Session.Token)
}

type Config struct {
	Host          string
	Port          string
	Tls           TlsConfig
	Clients       Clients
	SessionSecret string
	TokenSecret   string
}

type TlsConfig struct {
	Enabled bool
	Key     string
	Cert    string
}

type Clients struct {
	Members cclient.MemberClient
	User    uclient.UserClient
}
