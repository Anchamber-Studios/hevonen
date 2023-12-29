package main

import (
	"github.com/anchamber-studios/hevonen/frontend/pages"
	"github.com/anchamber-studios/hevonen/frontend/types"
	"github.com/labstack/echo/v4"
)

const (
	HX_REQUEST_HEADER = "hx-request"
)

func index(c echo.Context) error {
	cc := c.(*types.CustomContext)
	hxRequest := cc.Request().Header.Get(HX_REQUEST_HEADER)
	if hxRequest == "true" {
		return pages.Index().Render(c.Request().Context(), c.Response().Writer)
	}
	return pages.IndexWL(cc.Session).Render(c.Request().Context(), c.Response().Writer)
}
