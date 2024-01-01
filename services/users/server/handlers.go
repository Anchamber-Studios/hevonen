package server

import (
	"net/http"

	"github.com/anchamber-studios/hevonen/services/users/client"
	"github.com/labstack/echo/v4"
)

const (
	PathParamMemberId = "userId"
)

type MemberHandler struct{}

func list(c echo.Context) error {
	cc := c.(*CustomContext)
	users, err := cc.Services.Users.List(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	return c.JSON(http.StatusOK, &users)
}

func new(c echo.Context) error {
	cc := c.(*CustomContext)
	var newUser client.UserCreate
	if err := cc.Bind(&newUser); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	user, err := cc.Services.Users.Create(c.Request().Context(), newUser)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	return c.JSON(http.StatusOK, &user)
}

func login(c echo.Context) error {
	cc := c.(*CustomContext)
	var login client.UserLogin
	if err := cc.Bind(&login); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}

	user, err := cc.Services.Users.Login(c.Request().Context(), login)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	return c.JSON(http.StatusOK, &user)
}

func details(c echo.Context) error {
	cc := c.(*CustomContext)
	user, err := cc.Services.Users.Get(c.Request().Context(), c.Param(PathParamMemberId))
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	return c.JSON(http.StatusOK, &user)
}
