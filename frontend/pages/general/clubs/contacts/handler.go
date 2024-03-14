package contacts

import (
	"github.com/anchamber-studios/hevonen/frontend/types"
	ctypes "github.com/anchamber-studios/hevonen/services/club/shared/types"
	"github.com/labstack/echo/v4"
)

func GetListContacts(c echo.Context) error {
	cc := c.(*types.CustomContext)
	clubID := cc.Param("clubID")
	contacts, err := cc.Config.Clients.Contacts.List(cc.ClientContext(), clubID)
	if err != nil {
		cc.Logger().Errorf("Unable to get contacts: %v\n", err)
		contacts = []ctypes.Contact{}
	}
	props := ContactListProps{
		Contacts: contacts,
	}
	if cc.HXRequest {
		return List(cc.Tr, props).Render(cc.Request().Context(), cc.Response().Writer)
	}
	return ListWL(cc.Session, cc.Tr, props).Render(cc.Request().Context(), cc.Response().Writer)
}
