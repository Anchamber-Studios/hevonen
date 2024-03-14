package server

import (
	"net/http"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/anchamber-studios/hevonen/services/admin/users/client"
	"github.com/labstack/echo/v4"
)

const (
	PathParamContactId = "userID"
)

type ContactHandler struct{}

func list(c echo.Context) error {
	cc := c.(*CustomContext)
	users, err := cc.Services.Users.List(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	return c.JSON(http.StatusOK, &users)
}

func new(c echo.Context) error {
	cc := c.(*CustomContext)
	var newUser client.UserCreate
	if err := cc.Bind(&newUser); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	user, err := cc.Services.Users.Create(c.Request().Context(), newUser)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	return c.JSON(http.StatusOK, &user)
}

func login(c echo.Context) error {
	cc := c.(*CustomContext)
	var login client.UserLogin
	if err := cc.Bind(&login); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}

	user, err := cc.Services.Users.Login(c.Request().Context(), login)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	key, err := paseto.V4SymmetricKeyFromHex(cc.Config.TokenSecret)
	if err != nil {
		cc.Logger().Errorf("Unable to create token key: %v\n", err)
	}
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))
	token.SetString("user-id", user.ID)
	signed := token.V4Encrypt(key, nil)

	return c.JSON(http.StatusOK, &client.UserLoginResponse{
		Token: signed,
		Email: user.Email,
		ID:    user.ID,
	})
}

func logout(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func details(c echo.Context) error {
	cc := c.(*CustomContext)
	user, err := cc.Services.Users.Get(c.Request().Context(), c.Param(PathParamContactId))
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	return c.JSON(http.StatusOK, &user)
}
