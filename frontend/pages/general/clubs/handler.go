package clubs

import (
	"github.com/anchamber-studios/hevonen/frontend/types"
	"github.com/labstack/echo/v4"
)

func GetNoClubs(c echo.Context) error {
	cc := c.(*types.CustomContext)

	if cc.HXRequest {
		return NoClubs().
			Render(cc.Request().Context(), cc.Response().Writer)
	}
	return NoClubsWL(cc.Session).
		Render(cc.Request().Context(), cc.Response().Writer)
}
