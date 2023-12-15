package main

import (
	"github.com/anchamber-studios/hevonen/frontend/pages"
	m "github.com/anchamber-studios/hevonen/frontend/pages/members"
	"github.com/labstack/echo/v4"
)

func index(c echo.Context) error {
	return pages.Index().Render(c.Request().Context(), c.Response().Writer)
}

func memberList(c echo.Context) error {
	cc := c.(*CustomContext)
	props := m.MemberListProps{
		Members: []m.Member{},
	}
	members, err := cc.Config.Clients.Members.GetMembers()
	if err != nil {
		c.Logger().Errorf("Unable to get members: %v\n", err)
		m.MemberList(props).Render(c.Request().Context(), c.Response().Writer)
	}
	for _, member := range members {
		props.Members = append(props.Members, m.Member{
			ID:        member.ID,
			FirstName: member.FirstName,
			LastName:  member.LastName,
			Email:     member.Email,
			Phone:     member.Phone,
		})
	}
	return m.MemberList(props).Render(c.Request().Context(), c.Response().Writer)
}
