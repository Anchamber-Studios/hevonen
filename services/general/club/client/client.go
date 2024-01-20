package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/anchamber-studios/hevonen/lib"
)

type ClubClient interface {
	ListClubsForIdentity(ctx lib.ClientContext, identityID string) ([]ClubMember, error)
}

type ClubClientHttp struct {
	Url string
}

func (c *ClubClientHttp) ListClubsForIdentity(ctx lib.ClientContext, identityID string) ([]ClubMember, error) {
	resp, err := http.Get(c.Url + "/i/" + identityID + "/c")
	if err != nil {
		return nil, err
	}
	var clubs []ClubMember
	err = json.NewDecoder(resp.Body).Decode(&clubs)
	if err != nil {
		return nil, err
	}
	return clubs, nil
}

type MemberClient interface {
	GetMembers() ([]Member, error)
	CreateMember(member MemberCreate) (string, error)
}
type MemberClientHttp struct {
	Url string
}

func (c *MemberClientHttp) GetMembers() ([]Member, error) {
	resp, err := http.Get(c.Url + "/members")
	if err != nil {
		return nil, err
	}
	var members []Member
	err = json.NewDecoder(resp.Body).Decode(&members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (c *MemberClientHttp) CreateMember(member MemberCreate) (string, error) {
	memberJson, err := json.Marshal(member)
	if err != nil {
		return "", err
	}
	res, err := http.Post(c.Url+"/members", "application/json", bytes.NewReader(memberJson))
	if err != nil {
		return "", err
	}
	location := res.Header.Get("Location")
	return location, nil
}
