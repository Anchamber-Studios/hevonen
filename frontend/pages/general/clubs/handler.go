package clubs

import (
	"net/http"

	"github.com/anchamber-studios/hevonen/frontend/types"
	"github.com/anchamber-studios/hevonen/services/club/client"
	"github.com/labstack/echo/v4"
)

func GetNoClubs(c echo.Context) error {
	cc := c.(*types.CustomContext)

	if cc.HXRequest {
		return NoClubs(cc.Tr).
			Render(cc.Request().Context(), cc.Response().Writer)
	}
	return NoClubsWL(cc.Session, cc.Tr).
		Render(cc.Request().Context(), cc.Response().Writer)
}

func GetCreateForm(c echo.Context) error {
	cc := c.(*types.CustomContext)

	if cc.HXRequest {
		return CreateForm(cc.Tr, CreateClubFormProps{}).
			Render(cc.Request().Context(), cc.Response().Writer)
	}
	return CreateFormWL(cc.Session, cc.Tr, CreateClubFormProps{}).
		Render(cc.Request().Context(), cc.Response().Writer)
}

func PostCreateForm(c echo.Context) error {
	cc := c.(*types.CustomContext)

	var club client.ClubCreate
	if err := cc.Bind(&club); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}

	valErr := client.ValidateClubCreate(club)
	if len(valErr.Children) > 0 {
		cc.Response().Status = http.StatusBadRequest
		return CreateForm(cc.Tr, CreateClubFormProps{
			Errors: valErr.Children,
			Values: club,
		}).
			Render(cc.Request().Context(), cc.Response().Writer)
	}

	cId, err := cc.Config.Clients.Clubs.CreateClub(cc.ClientContext(), club)
	if err != nil {
		cc.Logger().Errorf("Unable to create club: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cc.Response().Header().Set("HX-Redirect", "/c/"+cId)
	return c.NoContent(http.StatusTemporaryRedirect)
}
