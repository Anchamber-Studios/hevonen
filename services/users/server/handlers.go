package server

import (
	"github.com/labstack/echo/v4"
)

const (
	PathParamMemberId = "memberId"
)

type MemberHandler struct{}

func list(c echo.Context) error {
	return echo.NewHTTPError(echo.ErrNotImplemented.Code, "not implemented")
}

func new(c echo.Context) error {
	return echo.NewHTTPError(echo.ErrNotImplemented.Code, "not implemented")
}

func details(c echo.Context) error {
	return echo.NewHTTPError(echo.ErrNotImplemented.Code, "not implemented")
}
