package auth

import (
	"net/http"

	"github.com/anchamber-studios/hevonen/frontend/translation"
	"github.com/anchamber-studios/hevonen/frontend/types"
	"github.com/anchamber-studios/hevonen/services/admin/users/client"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	ory "github.com/ory/client-go"
)

func GetLogin(c echo.Context) error {
	cc := c.(*types.CustomContext)
	csrf := c.Get("csrf").(string)
	// hxRequest := cc.Request().Header.Get(HX_REQUEST_HEADER)
	redirect := c.QueryParam("redirect")
	cc.Response().Header().Set("HX-Target", "html")
	localizer := i18n.NewLocalizer(translation.GetBundle(), language.English.String())
	return LoginPage(csrf, LoginPageProps{RedirectUrl: redirect}, localizer).Render(c.Request().Context(), c.Response().Writer)
}

func PostLogin(c echo.Context) error {
	cc := c.(*types.CustomContext)
	login := client.UserLogin{}
	if err := cc.Bind(&login); err != nil {
		cc.Logger().Errorf("Unable to bind login: %v\n", err)
		return c.String(http.StatusUnauthorized, "invalid logn")
	}

	body := ory.UpdateLoginFlowBody{
		UpdateLoginFlowWithPasswordMethod: &ory.UpdateLoginFlowWithPasswordMethod{
			Password:   login.Password,
			Identifier: login.Email,
			Method:     "password",
			AdditionalProperties: map[string]interface{}{
				"expires_at": "60",
			},
		},
	}
	flow, _, err := cc.Auth.FrontendAPI.CreateNativeLoginFlow(c.Request().Context()).
		Refresh(false).
		Execute()
	if err != nil {
		cc.Logger().Errorf("Unable to create login flow: %v\n", err)
		return err
	}
	flowResponse, _, err := cc.Auth.FrontendAPI.UpdateLoginFlow(c.Request().Context()).
		Flow(flow.Id).
		UpdateLoginFlowBody(body).
		Execute()
	if err != nil {
		cc.Logger().Errorf("Unable to login: %v\n", err)
		csrf := cc.Get("csrf").(string)
		localizer := i18n.NewLocalizer(translation.GetBundle(), language.English.String())
		return LoginForm(csrf, LoginPageProps{
			EmailError:    "",
			PasswordError: "",
			Error:         "email or password is incorrect",
		}, localizer).Render(cc.Request().Context(), cc.Response().Writer)
	}

	cc.Logger().Infof("user '%s' logged in\n", flowResponse.GetSession().Identity.Id)
	cc.Response().Header().Set("HX-Target", "html")
	rediectUrl := cc.FormValue("redirect")
	if rediectUrl != "" {
		cc.Response().Header().Set("HX-Redirect", rediectUrl)
	} else {
		cc.Response().Header().Set("HX-Redirect", "/")
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    *flowResponse.SessionToken,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(cc.Response().Writer, cookie)

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
