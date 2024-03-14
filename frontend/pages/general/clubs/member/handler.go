package member

import (
	"github.com/anchamber-studios/hevonen/frontend/types"
	ctypes "github.com/anchamber-studios/hevonen/services/club/shared/types"
	"github.com/labstack/echo/v4"
)

func GetListMembers(c echo.Context) error {
	cc := c.(*types.CustomContext)
	clubID := cc.Param("clubID")
	members, err := cc.Config.Clients.Members.List(cc.ClientContext(), clubID)
	if err != nil {
		cc.Logger().Errorf("Unable to get members: %v\n", err)
		members = []ctypes.Member{}
	}
	props := MemberListProps{
		Members: members,
	}
	if cc.HXRequest {
		return List(cc.Tr, props).Render(cc.Request().Context(), cc.Response().Writer)
	}
	return ListWL(cc.Session, cc.Tr, props).Render(cc.Request().Context(), cc.Response().Writer)
}
