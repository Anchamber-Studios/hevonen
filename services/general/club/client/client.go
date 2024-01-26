package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anchamber-studios/hevonen/lib"
)

type ClubClient interface {
	ListClubsForIdentity(ctx lib.ClientContext) ([]ClubMember, error)
	CreateClub(ctx lib.ClientContext, club ClubCreate) (string, error)
}

type ClubClientHttp struct {
	Url string
}

func (c *ClubClientHttp) ListClubsForIdentity(ctx lib.ClientContext) ([]ClubMember, error) {
	resp, err := http.Get(c.Url + "/i/" + ctx.IdentityID + "/c")
	if err != nil {
		return nil, err
	}
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return nil, fmt.Errorf("unable to get clubs for identity (%d): %v", resp.StatusCode, resp.Body)
	}
	var clubs []ClubMember
	err = json.NewDecoder(resp.Body).Decode(&clubs)
	if err != nil {
		return nil, err
	}
	return clubs, nil
}

func (c *ClubClientHttp) CreateClub(ctx lib.ClientContext, club ClubCreate) (string, error) {
	valErr := ValidateClubCreate(club)
	if len(valErr.Children) > 0 {
		return "", valErr
	}

	clubJson, err := json.Marshal(club)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(c.Url+"/i/"+ctx.IdentityID+"/c", "application/json", bytes.NewReader(clubJson))
	if err != nil {
		return "", err
	}
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return "", fmt.Errorf("unable to get clubs for identity (%d): %v", resp.StatusCode, resp.Body)
	}
	location := resp.Header.Get("Location")
	return location, nil
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
