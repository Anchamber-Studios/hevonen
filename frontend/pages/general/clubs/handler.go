package clubs

import (
	"net/http"

	"github.com/anchamber-studios/hevonen/frontend/types"
	ctypes "github.com/anchamber-studios/hevonen/services/club/shared/types"
	"github.com/labstack/echo/v4"
)

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

	var club ctypes.ClubCreate
	if err := cc.Bind(&club); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}

	valErr := ctypes.ValidateClubCreate(club)
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
		header := types.HxTriggerHeader{ErrorMessage: "Unable to create club"}
		cc.Response().Header().Add(header.Key(), header.ToJSON())
		return cc.NoContent(http.StatusNoContent)
		// return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cc.Response().Header().Set("HX-Redirect", "/clubs/"+cId)
	return c.NoContent(http.StatusTemporaryRedirect)
}

func GetListClubs(c echo.Context) error {
	cc := c.(*types.CustomContext)
	userClubs, err := cc.Config.Clients.Clubs.List(cc.ClientContext())
	if err != nil {
		cc.Logger().Errorf("Unable to get clubs: %v\n", err)
		return cc.Redirect(302, "/clubs")
	}
	if cc.HXRequest {
		return List(cc.Tr, ListProps{
			Clubs: userClubs,
		}).
			Render(cc.Request().Context(), cc.Response().Writer)
	}
	return ListWL(cc.Session, cc.Tr, ListProps{
		Clubs: userClubs,
	}).Render(cc.Request().Context(), cc.Response().Writer)
}
