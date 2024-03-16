package calendar

import (
	"strconv"
	"time"

	"github.com/anchamber-studios/hevonen/frontend/types"
	"github.com/labstack/echo/v4"
)

func GetDaySelection(c echo.Context) error {
	cc := c.(*types.CustomContext)

	year, err := strconv.ParseInt(cc.QueryParam("year"), 10, 32)
	if err != nil || year < 2020 || year > 2500 {
		cc.Logger().Errorf("Invalid year: %v\n", err)
		year = int64(time.Now().Year())
	}
	month, err := strconv.ParseInt(cc.QueryParam("month"), 10, 32)
	if err != nil || month < 1 || month > 12 {
		cc.Logger().Errorf("Invalid month: %v\n", err)
		month = int64(time.Now().Month())
	}
	return DaySelection(cc.Tr, DaySelectionProps{
		NumberOfMonth: 2,
		Year:          int(year),
		Month:         time.Month(month),
	}).Render(cc.Request().Context(), cc.Response().Writer)
}
