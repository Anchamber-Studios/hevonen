package server

import (
	"fmt"
	"net/http"

	"github.com/anchamber-studios/hevonen/lib"
	"github.com/anchamber-studios/hevonen/services/club/shared/types"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

const (
	PathParamContactId = "memberId"
	PathIdentityId     = "identityId"
)

type ClubHandler struct{}

// Lists all clubs the current user is a contact of
func (h *ClubHandler) List(c echo.Context) error {
	cc := c.(*CustomContext)
	identityID := cc.Get("identityID").(string)
	clubs, err := cc.Services.Clubs.List(c.Request().Context())
	if err != nil {
		return handleServiceError(cc, err)
	}
	cc.Logger().Infof("Found %d clubs for identity %s\n", len(clubs), identityID)
	return c.JSON(http.StatusOK, &clubs)
}

// Lists all clubs the current user is a contact of
func (h *ClubHandler) ListForIdentity(c echo.Context) error {
	cc := c.(*CustomContext)
	identityID := cc.Get("identityID").(string)
	clubs, err := cc.Services.Clubs.ListForIdentity(c.Request().Context(), identityID)
	if err != nil {
		return handleServiceError(cc, err)
	}
	cc.Logger().Infof("Found %d clubs for identity %s\n", len(clubs), identityID)
	return c.JSON(http.StatusOK, &clubs)
}

// Create a new club. The current user will be added as admin contact of the club
func (h *ClubHandler) Create(c echo.Context) error {
	cc := c.(*CustomContext)
	identityID := cc.Get("identityID").(string)
	var club types.ClubCreate
	if err := c.Bind(&club); err != nil {
		cc.Logger().Errorf("Unable to bind club: %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to bind club")
	}

	valErr := types.ValidateClubCreate(club)
	if len(valErr.Children) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, valErr)
	}

	email := cc.Get("email").(string)
	contact := types.ContactCreate{
		IdentityID: identityID,
		Email:      email,
	}
	cId, err := cc.Services.Clubs.CreateWithAdminContact(c.Request().Context(), club, contact)
	if err != nil {
		return handleServiceError(cc, err)
	}
	cc.Logger().Infof("Created club %s\n", cId)
	cc.Response().Header().Set("Location", fmt.Sprintf("/clubs/%s", cId))
	return cc.NoContent(http.StatusCreated)
}

func (h *ClubHandler) DeleteClub(c echo.Context) error {
	cc := c.(*CustomContext)
	identityID := cc.Get("identityID").(string)
	clubID := c.Param("clubID")
	err := cc.Services.Clubs.Delete(c.Request().Context(), identityID, clubID)
	if err != nil {
		return handleServiceError(cc, err)
	}
	cc.Logger().Infof("Deleted club %s\n", clubID)
	return cc.NoContent(http.StatusNoContent)
}

type ContactHandler struct{}

func (h *ContactHandler) List(c echo.Context) error {
	cc := c.(*CustomContext)
	cID := c.Param("clubID")
	contacts, err := cc.Services.Contacts.List(cc.Request().Context(), cID)
	if err != nil {
		return handleServiceError(cc, err)
	}
	cc.Logger().Infof("Found %d contacts for club %s\n", len(contacts), cID)
	return c.JSON(http.StatusOK, &contacts)
}

func handleServiceError(_ *CustomContext, err error) error {
	if e, ok := err.(*pgconn.PgError); ok {
		if e.Code == "23505" {
			return echo.NewHTTPError(http.StatusConflict, lib.NewAlreadyExistsError("club", "name", "Test Club"))
		}
	}
	return echo.NewHTTPError(http.StatusInternalServerError, lib.NewInternalServerError())
}
