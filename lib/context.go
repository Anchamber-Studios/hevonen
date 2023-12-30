package lib

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ClientContext struct {
	Token     string
	RequestID string
}

func NewClientContext(c echo.Context, token string) ClientContext {
	return ClientContext{
		Token:     token,
		RequestID: c.Request().Header.Get(echo.HeaderXRequestID),
	}
}

func (cc ClientContext) SetHeader(req *http.Request) {
	req.Header.Set(echo.HeaderXRequestID, cc.RequestID)
	if cc.Token != "" {
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", cc.Token))
	}
}
