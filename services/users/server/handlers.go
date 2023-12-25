package server

import (
	"net/http"

	"github.com/anchamber-studios/hevonen/services/users/client"
	"github.com/labstack/echo/v4"
)

const (
	PathParamMemberId = "memberId"
)

type MemberHandler struct{}

func list(c echo.Context) error {
	return echo.NewHTTPError(echo.ErrNotImplemented.Code, "not implemented")
}

func new(c echo.Context) error {
	return echo.NewHTTPError(echo.ErrNotImplemented.Code, "not implemented")
}

func login(c echo.Context) error {
	cc := c.(*CustomContext)
	var login client.UserLogin
	if err := cc.Bind(&login); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	user, err := cc.Repos.Users.Login(c.Request().Context(), login)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	return c.JSON(http.StatusOK, &user)
}

func details(c echo.Context) error {
	return echo.NewHTTPError(echo.ErrNotImplemented.Code, "not implemented")
}