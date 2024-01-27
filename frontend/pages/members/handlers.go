package members

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/anchamber-studios/hevonen/frontend/types"
	ctypes "github.com/anchamber-studios/hevonen/services/club/shared/types"
	"github.com/labstack/echo/v4"
)

func GetMemberList(c echo.Context) error {
	cc := c.(*types.CustomContext)
	props := MemberListProps{
		Members: []ctypes.Member{},
	}
	members, err := cc.Config.Clients.Members.GetMembers()
	if err != nil {
		c.Logger().Errorf("Unable to get members: %v\n", err)
		members = []ctypes.Member{}
	}
	for _, member := range members {
		props.Members = append(props.Members, ctypes.Member{
			ID:         member.ID,
			FirstName:  member.FirstName,
			MiddleName: member.MiddleName,
			LastName:   member.LastName,
			Email:      member.Email,
			Phone:      member.Phone,
		})
	}
	if cc.HXRequest {
		return MemberList(props).Render(c.Request().Context(), c.Response().Writer)
	}
	cc.Logger().Errorf("render with layout\n")
	return MemberListWL(cc.Session, props).Render(c.Request().Context(), c.Response().Writer)
}

func GetNewMemberForm(c echo.Context) error {
	cc := c.(*types.CustomContext)
	csrf := c.Get("csrf").(string)
	var tmpl templ.Component
	if cc.HXRequest {
		tmpl = NewForm(MemberFormProps{Csrf: csrf})
	} else {
		tmpl = NewFormWL(cc.Session, MemberFormProps{Csrf: csrf})
	}
	return tmpl.Render(c.Request().Context(), c.Response().Writer)
}

func PostNewMemberForm(c echo.Context) error {
	cc := c.(*types.CustomContext)
	member := ctypes.MemberCreate{}
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
