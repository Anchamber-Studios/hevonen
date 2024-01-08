package server

import (
	"net/http"

	"github.com/anchamber-studios/hevonen/services/general/profile/client"
	"github.com/labstack/echo/v4"
)

const (
	PathParamProfileId = "profileId"
)

type MemberHandler struct{}

func new(c echo.Context) error {
	cc := c.(*CustomContext)
	var newProfile client.ProfileCreateRequest
	if err := cc.Bind(&newProfile); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	profile, err := cc.Services.Profiles.Create(c.Request().Context(), newProfile)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	return c.JSON(http.StatusOK, &profile)
}

func details(c echo.Context) error {
	cc := c.(*CustomContext)
	profile, err := cc.Services.Profiles.Get(c.Request().Context(), c.Param(PathParamProfileId))
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	return c.JSON(http.StatusOK, &profile)
}
