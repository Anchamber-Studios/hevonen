package lib

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ClientContext struct {
	Token      string
	RequestID  string
	IdentityID string
	Email      string
}

func NewClientContext(c echo.Context, token string, identityID string, email string) ClientContext {
	return ClientContext{
		Token:      token,
		RequestID:  c.Request().Header.Get(echo.HeaderXRequestID),
		IdentityID: identityID,
		Email:      email,
	}
}

func (cc ClientContext) SetHeader(req *http.Request) {
	req.Header.Set(echo.HeaderXRequestID, cc.RequestID)
	if cc.Token != "" {
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", cc.Token))
		req.Header.Set("X-Request-ID", cc.RequestID)
	}
}
