package admin

import (
	"github.com/anchamber-studios/hevonen/frontend/types"
	"github.com/labstack/echo/v4"
)

func GetUsers(e echo.Context) error {
	cc := e.(*types.CustomContext)
	users, err := cc.Config.Clients.User.GetUsers(cc.ClientContext())
	if err != nil {
		return err
	}
	if cc.HXRequest {
		return UserList(UserListProps{Users: users}).Render(e.Request().Context(), e.Response().Writer)
	}
	return UserListWL(cc.Session, UserListProps{Users: users}).Render(e.Request().Context(), e.Response().Writer)
}

func GetUser(e echo.Context) error {
	cc := e.(*types.CustomContext)
	user, err := cc.Config.Clients.User.GetUser(cc.ClientContext(), e.Param("userID"))
	if err != nil {
		e.Logger().Errorf("Unable to get user: %v\n", err)
		return err
	}
	if cc.HXRequest {
		return UserDetails(UserDetailsProps{User: user}).Render(e.Request().Context(), e.Response().Writer)
	}
	return UserDetailsWL(cc.Session, UserDetailsProps{User: user}).Render(e.Request().Context(), e.Response().Writer)
}
