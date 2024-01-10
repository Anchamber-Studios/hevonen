package server

import (
	"net/http"

	"github.com/anchamber-studios/hevonen/lib"
	"github.com/anchamber-studios/hevonen/services/general/profile/client"
	"github.com/labstack/echo/v4"
)

const (
	PathParamProfileId = "profileIdnetityId"
)

type MemberHandler struct{}

func new(c echo.Context) error {
	cc := c.(*CustomContext)
	var newProfile client.ProfileCreateRequest
	if err := cc.Bind(&newProfile); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	profileId, err := cc.Services.Profiles.Create(c.Request().Context(), newProfile)
	if err != nil {
		cc.Logger().Errorf("Unable to create profile: %v\n", err)
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	cc.Logger().Infof("Creted profile %s\n", profileId)
	cc.Response().Header().Set("Location", "/profiles/"+profileId)
	return cc.NoContent(http.StatusNoContent)
}

func details(c echo.Context) error {
	cc := c.(*CustomContext)
	profileId := c.Param(PathParamProfileId)
	profile, err := cc.Services.Profiles.GetByIdentityID(c.Request().Context(), profileId)
	if err != nil {
		cc.Logger().Errorf("Unable to get profile for %s: %v\n", profileId, err)
		if e, ok := err.(*lib.ApiError); ok {
			return c.JSON(e.StatusCode, e)
		}
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	cc.Logger().Infof("Found profile %v\n", profileId)
	return cc.JSON(http.StatusOK, &profile)
}

func update(c echo.Context) error {
	cc := c.(*CustomContext)
	profileId := c.Param(PathParamProfileId)
	cc.Logger().Infof("Updating profile %v\n", profileId)
	var updateProfile client.ProfileUpdateRequest
	if err := cc.Bind(&updateProfile); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	err := cc.Services.Profiles.Update(c.Request().Context(), profileId, updateProfile)
	if err != nil {
		cc.Logger().Errorf("Unable to update profile for %s: %v\n", profileId, err)
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	cc.Logger().Infof("Updated profile %v\n", profileId)
	return cc.NoContent(http.StatusNoContent)
}
