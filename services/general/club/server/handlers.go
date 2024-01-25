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
	PathIdentityId    = "identityId"
)

type ClubHandler struct{}

func (h *ClubHandler) ListForIdentity(c echo.Context) error {
	cc := c.(*CustomContext)
	identityId := c.Param(PathIdentityId)
	clubs, err := cc.Repos.Clubs.ListForIdentity(c.Request().Context(), identityId)
	if err != nil {
		cc.Logger().Errorf("Unable to get clubs for identity %s: %v\n", identityId, err)
		echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	cc.Logger().Infof("Found %d clubs for identity %s\n", len(clubs), identityId)
	return c.JSON(http.StatusOK, &clubs)
}

func (h *ClubHandler) Create(c echo.Context) error {
	cc := c.(*CustomContext)
	identityId := c.Param(PathIdentityId)
	var club client.ClubCreate
	if err := c.Bind(&club); err != nil {
		cc.Logger().Errorf("Unable to bind club: %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to bind club")
	}

	valErr := client.ValidateClubCreate(club)
	if len(valErr.Children) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, valErr)
	}

	cId, err := cc.Repos.Clubs.Create(c.Request().Context(), club)
	if err != nil {
		cc.Logger().Errorf("Unable to create club for identity %s: %v\n", identityId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	cc.Response().Header().Set("Location", fmt.Sprintf("/i/%s/c/%s", identityId, cId))
	return cc.NoContent(http.StatusCreated)
}

type MemberHandler struct{}

func (h *MemberHandler) list(c echo.Context) error {
	cc := c.(*CustomContext)
	members, err := cc.Repos.Members.List(c.Request().Context())
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, &members)
}

func (h *MemberHandler) new(c echo.Context) error {
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

func (h *MemberHandler) details(c echo.Context) error {
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
