package members

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/anchamber-studios/hevonen/frontend/types"
	"github.com/anchamber-studios/hevonen/services/club/client"
	"github.com/labstack/echo/v4"
)

const (
	HX_REQUEST_HEADER = "hx-request"
)

func GetMemberList(c echo.Context) error {
	cc := c.(*types.CustomContext)
	props := MemberListProps{
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
		return MemberList(props).Render(c.Request().Context(), c.Response().Writer)
	}
	cc.Logger().Errorf("render with layout\n")
	return MemberListWL(cc.Session, props).Render(c.Request().Context(), c.Response().Writer)
}

func GetNewMemberForm(c echo.Context) error {
	cc := c.(*types.CustomContext)
	csrf := c.Get("csrf").(string)
	hxRequest := cc.Request().Header.Get(HX_REQUEST_HEADER)
	var tmpl templ.Component
	if hxRequest == "true" {
		tmpl = NewForm(MemberFormProps{Csrf: csrf})
	} else {
		tmpl = NewFormWL(cc.Session, MemberFormProps{Csrf: csrf})
	}
	return tmpl.Render(c.Request().Context(), c.Response().Writer)
}

func PostNewMemberForm(c echo.Context) error {
	cc := c.(*types.CustomContext)
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
	return c.Redirect(http.StatusFound, loc)
}
