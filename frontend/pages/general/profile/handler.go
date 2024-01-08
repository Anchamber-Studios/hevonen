package profile

import (
	"github.com/anchamber-studios/hevonen/frontend/types"
	"github.com/labstack/echo/v4"
)

func GetProfile(e echo.Context) error {
	cc := e.(*types.CustomContext)
	user, err := cc.Config.Clients.Profile.GetProfile(cc.ClientContext(), cc.Session.ID)
	if err != nil {
		e.Logger().Errorf("Unable to get user: %v\n", err)
		return err
	}
	props := ProfilePageProps{
		Profile: &Profile{
			FirstName:  user.FirstName,
			MiddleName: user.MiddleName,
			LastName:   user.LastName,
			Height:     user.Height,
			Weight:     user.Weight,
			Birthday:   user.Birthday,
		},
	}
	if cc.HXRequest {
		return ProfilePage(props).
			Render(e.Request().Context(), e.Response().Writer)
	}
	return ProfilePageWL(cc.Session, props).
		Render(e.Request().Context(), e.Response().Writer)
}
