package server

import (
	"github.com/labstack/echo/v4"
)

func GetGroups(e echo.Context) error {
	cc := e.(*CustomContext)
	groups, err := cc.Services.AuthService.GetGroups(cc.Request().Context())
	if err != nil {
		return echo.NewHTTPError(500, "Could not get groups")
	}
	return e.JSON(200, groups)
}

func GetAuthorizations(e echo.Context) error {
	cc := e.(*CustomContext)
	authorizations, err := cc.Services.AuthService.GetAuthorizations(cc.Request().Context())
	if err != nil {
		return echo.NewHTTPError(500, "Could not get authorizations")
	}
	return e.JSON(200, authorizations)
}

func GetAuthorizationsForService(e echo.Context) error {
	cc := e.(*CustomContext)
	serviceId := cc.Param("serviceId")
	authorizations, err := cc.Services.AuthService.GetAuthorizationsForService(cc.Request().Context(), serviceId)
	if err != nil {
		return echo.NewHTTPError(500, "Could not get authorizations")
	}
	return e.JSON(200, authorizations)
}

func RegisterServiceAuthorizations(e echo.Context) error {
	return echo.NewHTTPError(500, "Not implemented")
}

func DeleteServiceAuthorizations(e echo.Context) error {
	return echo.NewHTTPError(500, "Not implemented")
}

func UpdateServiceAuthorizations(e echo.Context) error {
	return echo.NewHTTPError(500, "Not implemented")
}

func GetServices(e echo.Context) error {
	cc := e.(*CustomContext)
	services, err := cc.Services.AuthService.GetServices(cc.Request().Context())
	if err != nil {
		return echo.NewHTTPError(500, "Could not get services")
	}
	return e.JSON(200, services)
}

func RegisterSercice(e echo.Context) error {
	return echo.NewHTTPError(500, "Not implemented")
}

func DeleteService(e echo.Context) error {
	return echo.NewHTTPError(500, "Not implemented")
}
