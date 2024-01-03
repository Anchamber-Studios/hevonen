package auth

import (
	"net/http"

	"github.com/anchamber-studios/hevonen/frontend/types"
	"github.com/anchamber-studios/hevonen/lib"
	"github.com/anchamber-studios/hevonen/services/admin/users/client"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func GetLogin(c echo.Context) error {
	cc := c.(*types.CustomContext)
	csrf := c.Get("csrf").(string)
	// hxRequest := cc.Request().Header.Get(HX_REQUEST_HEADER)
	cc.Response().Header().Set("HX-Target", "html")
	return LoginPage(csrf, LoginPageProps{}).Render(c.Request().Context(), c.Response().Writer)
}

func PostLogin(c echo.Context) error {
	cc := c.(*types.CustomContext)
	login := client.UserLogin{}
	if err := cc.Bind(&login); err != nil {
		cc.Logger().Errorf("Unable to bind login: %v\n", err)
		return c.String(http.StatusUnauthorized, "invalid logn")
	}

	user, err := cc.Config.Clients.User.Login(lib.NewClientContext(c, ""), login)
	if err != nil {
		cc.Logger().Errorf("Unable to login: %v\n", err)
		csrf := cc.Get("csrf").(string)
		return LoginForm(csrf, LoginPageProps{
			EmailError:    "",
			PasswordError: "",
			Error:         "email or password is incorrect",
		}).Render(cc.Request().Context(), cc.Response().Writer)
	}
	cc.Logger().Errorf("user '%s' logged in\n", user.ID)
	cc.Response().Header().Set("HX-Target", "html")

	cc.Response().Header().Set("HX-Redirect", "/")
	sess, _ := session.Get("session", cc)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 2,
		HttpOnly: true,
	}
	sess.Values["id"] = user.ID
	sess.Values["email"] = user.Email
	sess.Values["token"] = user.Token
	cc.Logger().Errorf("save session")
	if err := sess.Save(cc.Request(), cc.Response()); err != nil {
		cc.Logger().Errorf("Unable to save session: %v\n", err)
	}
	return c.NoContent(http.StatusNoContent)
}

func GetRegister(c echo.Context) error {
	cc := c.(*types.CustomContext)
	csrf := c.Get("csrf").(string)
	// hxRequest := cc.Request().Header.Get(HX_REQUEST_HEADER)
	cc.Response().Header().Set("HX-Target", "html")
	return RegisterPage(csrf, RegisterPageProps{}).Render(c.Request().Context(), c.Response().Writer)
}

func GetLogout(c echo.Context) error {
	cc := c.(*types.CustomContext)
	cc.Response().Header().Set("HX-Target", "html")
	sess, _ := session.Get("session", cc)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	if err := sess.Save(cc.Request(), cc.Response()); err != nil {
		cc.Logger().Errorf("Unable to save session: %v\n", err)
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
}
