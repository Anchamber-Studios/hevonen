package profile

import (
	"fmt"
	"time"

	"github.com/anchamber-studios/hevonen/frontend/types"
	"github.com/anchamber-studios/hevonen/lib"
	"github.com/anchamber-studios/hevonen/services/general/profile/client"
	"github.com/labstack/echo/v4"
)

func GetProfile(c echo.Context) error {
	cc := c.(*types.CustomContext)
	profile, err := cc.Config.Clients.Profile.GetProfile(cc.ClientContext(), cc.Session.ID)
	if err != nil {
		cc.Logger().Errorf("profile for identity not found, creating new profile")
		if apiErr, ok := err.(*lib.ApiError); ok {
			cc.Logger().Errorf("profile for identity not found, creating new profile")
			_, err := createIfNotExists(cc, apiErr)
			if err != nil {
				cc.Logger().Errorf("Unable to create profile: %v\n", err)
				return err
			}
			profile = client.ProfileResponse{}
		} else {
			cc.Logger().Errorf("Unable to get user: %v\n", err)
			return err
		}
	}

	props := ProfilePageProps{
		Profile: &Profile{
			FirstName:  profile.FirstName,
			MiddleName: profile.MiddleName,
			LastName:   profile.LastName,
			Height:     profile.Height,
			Weight:     profile.Weight,
			Birthday:   profile.Birthday.Format(time.DateOnly),
		},
	}
	if cc.HXRequest {
		return ProfilePage(props).
			Render(cc.Request().Context(), cc.Response().Writer)
	}
	return ProfilePageWL(cc.Session, props).
		Render(cc.Request().Context(), cc.Response().Writer)
}

func UpdateProfile(c echo.Context) error {
	cc := c.(*types.CustomContext)
	var profile Profile
	if err := cc.Bind(&profile); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	fmt.Printf("profile: %v\n", profile)
	birthday, err := time.Parse(time.DateOnly, profile.Birthday)
	err = cc.Config.Clients.Profile.UpdateProfile(cc.ClientContext(), cc.Session.ID, client.ProfileUpdateRequest{
		FirstName:  profile.FirstName,
		MiddleName: profile.MiddleName,
		LastName:   profile.LastName,
		Height:     profile.Height,
		Weight:     profile.Weight,
		Birthday:   birthday,
	})
	if err != nil {
		cc.Logger().Errorf("Unable to update profile for %s: %v\n", cc.Session.ID, err)
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	cc.Logger().Infof("Updated profile %v\n", cc.Session.ID)
	cc.Logger().Infof("%v\n", profile)
	props := ProfilePageProps{
		Profile: &profile,
	}
	return profileForm(props).
		Render(cc.Request().Context(), cc.Response().Writer)
}

func createIfNotExists(cc *types.CustomContext, err *lib.ApiError) (string, error) {
	if err.StatusCode == 404 {
		return cc.Config.Clients.Profile.CreateProfile(cc.ClientContext(), client.ProfileCreateRequest{
			IdentityID: cc.Session.ID,
		})
	}
	return "", err
}
