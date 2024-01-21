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
	cc.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusTemporaryRedirect)
}
