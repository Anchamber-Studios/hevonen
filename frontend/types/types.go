package types

import (
	"time"

	"github.com/anchamber-studios/hevonen/lib"
	uclient "github.com/anchamber-studios/hevonen/services/admin/users/client"
	cclient "github.com/anchamber-studios/hevonen/services/club/client"
	pclient "github.com/anchamber-studios/hevonen/services/general/profile/client"
	"github.com/labstack/echo/v4"
	ory "github.com/ory/client-go"
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
	Auth      *ory.APIClient
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
	Profile pclient.ProfileClient
}
type FormDate time.Time

func (ct *FormDate) UnmarshalParam(param string) error {
	t, err := time.Parse(`2006-01-02 15:04:05 MST`, param)
	if err != nil {
		return err
	}
	*ct = FormDate(t)
	return nil
}
