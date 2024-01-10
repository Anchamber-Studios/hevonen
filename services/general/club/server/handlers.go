package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/anchamber-studios/hevonen/lib"
	"github.com/anchamber-studios/hevonen/services/club/client"
	"github.com/labstack/echo/v4"
)

const (
	PathParamMemberId = "memberId"
)

type MemberHandler struct{}

func list(c echo.Context) error {
	cc := c.(*CustomContext)
	members, err := cc.Repos.Members.List(c.Request().Context())
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, &members)
}

func new(c echo.Context) error {
	cc := c.(*CustomContext)
	var member client.MemberCreate
	if err := c.Bind(&member); err != nil {
		c.Logger().Errorf("Unable to bind member: %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to bind member")
	}
	cId, err := cc.Repos.Members.Create(c.Request().Context(), member)
	if err != nil {
		c.Logger().Errorf("Unable to create member: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	cc.Response().Header().Set("Location", fmt.Sprintf("/members/%s", cId))
	return cc.NoContent(http.StatusCreated)
}

func details(c echo.Context) error {
	cc := c.(*CustomContext)
	memberIdEncoded := cc.Param(PathParamMemberId)
	member, err := cc.Repos.Members.Get(c.Request().Context(), memberIdEncoded)
	if err != nil {
		cc.Logger().Errorf("Unable to get member: %v\n", err)
		if errors.Is(err, lib.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "member not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, &member)
}
