package main

import (
	"github.com/anchamber-studios/hevonen/frontend/pages"
	a "github.com/anchamber-studios/hevonen/frontend/pages/auth"
	m "github.com/anchamber-studios/hevonen/frontend/pages/members"
	"github.com/anchamber-studios/hevonen/services/club/client"
	"github.com/labstack/echo/v4"
)

const (
	HX_REQUEST_HEADER = "hx-request"
)

func index(c echo.Context) error {
	cc := c.(*CustomContext)
	hxRequest := cc.Request().Header.Get(HX_REQUEST_HEADER)
	if hxRequest == "true" {
		return pages.Index().Render(c.Request().Context(), c.Response().Writer)
	}
	return pages.IndexWL().Render(c.Request().Context(), c.Response().Writer)
}

func memberList(c echo.Context) error {
	cc := c.(*CustomContext)
	props := m.MemberListProps{
		Members: []client.Member{},
	}
	members, err := cc.Config.Clients.Members.GetMembers()
	hxRequest := cc.Request().Header.Get(HX_REQUEST_HEADER)
	if err != nil {
		c.Logger().Errorf("Unable to get members: %v\n", err)
		members = []client.Member{}
	}
	for _, member := range members {
		props.Members = append(props.Members, client.Member{
			ID:         member.ID,
			FirstName:  member.FirstName,
			MiddleName: member.MiddleName,
			LastName:   member.LastName,
			Email:      member.Email,
			Phone:      member.Phone,
		})
	}
	cc.Logger().Errorf("HX Request: %v\n", hxRequest)
	if hxRequest == "true" {
		return m.MemberList(props).Render(c.Request().Context(), c.Response().Writer)
	}
	cc.Logger().Errorf("render with layout\n")
	return m.MemberListWL(props).Render(c.Request().Context(), c.Response().Writer)
}

func memberNew(c echo.Context) error {
	cc := c.(*CustomContext)
	csrf := c.Get("csrf").(string)
	hxRequest := cc.Request().Header.Get(HX_REQUEST_HEADER)
	if hxRequest == "true" {
		return m.NewForm(csrf, m.MemberFormProps{}).Render(c.Request().Context(), c.Response().Writer)
	}
	return m.NewFormWL(csrf, m.MemberFormProps{}).Render(c.Request().Context(), c.Response().Writer)
}

func postNewMember(c echo.Context) error {
	cc := c.(*CustomContext)
	member := client.MemberCreate{}
	if err := c.Bind(&member); err != nil {
		c.Logger().Errorf("Unable to bind member: %v\n", err)
		return c.String(500, "Unable to bind member")
	}
	c.Logger().Errorf("POST /members: MemberCreate%v\n", member)
	loc, err := cc.Config.Clients.Members.CreateMember(member)
	if err != nil {
		c.Logger().Errorf("Unable to add member: %v\n", err)
		return c.String(500, "Unable to add member")
	}
	return c.Redirect(302, loc)
}

func getLogin(c echo.Context) error {
	cc := c.(*CustomContext)
	csrf := c.Get("csrf").(string)
	// hxRequest := cc.Request().Header.Get(HX_REQUEST_HEADER)
	cc.Response().Header().Set("HX-Target", "html")
	return a.LoginPage(csrf, a.LoginPageProps{}).Render(c.Request().Context(), c.Response().Writer)
}

func getRegister(c echo.Context) error {
	cc := c.(*CustomContext)
	csrf := c.Get("csrf").(string)
	// hxRequest := cc.Request().Header.Get(HX_REQUEST_HEADER)
	cc.Response().Header().Set("HX-Target", "html")
	return a.RegisterPage(csrf, a.RegisterPageProps{}).Render(c.Request().Context(), c.Response().Writer)
}
