package calendar

import (
	"github.com/anchamber-studios/hevonen/frontend/types"
	"github.com/labstack/echo/v4"
)

func GetListCalendar(c echo.Context) error {
	cc := c.(*types.CustomContext)
	props := CalendarListProps{}
	if cc.HXRequest {
		return List(cc.Tr, props).
			Render(cc.Request().Context(), cc.Response().Writer)
	}
	return ListWL(cc.Session, cc.Tr, props).
		Render(cc.Request().Context(), cc.Response().Writer)
}
