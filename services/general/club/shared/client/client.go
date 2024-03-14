package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anchamber-studios/hevonen/lib"
	"github.com/anchamber-studios/hevonen/services/club/services"
	"github.com/anchamber-studios/hevonen/services/club/shared/types"
)

type ClubClient interface {
	ListClubsForIdentity(ctx lib.ClientContext) ([]types.ClubMember, error)
	List(ctx lib.ClientContext) ([]types.Club, error)
	CreateClub(ctx lib.ClientContext, club types.ClubCreate) (string, error)
}

type ClubClientHttp struct {
	Url     string
	Headers map[string]string
}

func (c *ClubClientHttp) WithHeader(key, value string) {
	if c.Headers == nil {
		c.Headers = make(map[string]string)
	}
	c.Headers[key] = value
}

func (c *ClubClientHttp) ListClubsForIdentity(ctx lib.ClientContext) ([]types.ClubMember, error) {
	req, err := http.NewRequest(http.MethodGet, c.Url+"/clubs/i", nil)
	if err != nil {
		return nil, err
	}
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	ctx.SetHeader(req)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return nil, fmt.Errorf("unable to get clubs for identity (%d): %v", resp.StatusCode, resp.Body)
	}
	var clubs []types.ClubMember
	err = json.NewDecoder(resp.Body).Decode(&clubs)
	if err != nil {
		return nil, err
	}
	return clubs, nil
}

func (c *ClubClientHttp) CreateClub(ctx lib.ClientContext, club types.ClubCreate) (string, error) {
	valErr := types.ValidateClubCreate(club)
	if len(valErr.Children) > 0 {
		return "", valErr
	}
	clubJson, err := json.Marshal(club)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, c.Url+"/clubs", bytes.NewReader(clubJson))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return "", err
	}
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	ctx.SetHeader(req)
	client := http.Client{}
	resp, err := client.Do(req)
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

func (c *ClubClientHttp) List(ctx lib.ClientContext) ([]types.Club, error) {
	req, err := http.NewRequest(http.MethodGet, c.Url+"/clubs", nil)
	if err != nil {
		return nil, err
	}
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	ctx.SetHeader(req)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return nil, fmt.Errorf("unable to get clubs (%d): %v", resp.StatusCode, resp.Body)
	}
	var clubs []types.Club
	err = json.NewDecoder(resp.Body).Decode(&clubs)
	if err != nil {
		return nil, err
	}
	return clubs, nil
}

type ClubClientLocal struct {
	clubService *services.ClubService
}

func (c *ClubClientLocal) ListClubsForIdentity(ctx lib.ClientContext) ([]types.ClubMember, error) {
	clubs, err := c.clubService.ListForIdentity(ctx.Context, ctx.IdentityID)
	if err != nil {
		return nil, err
	}
	return clubs, nil
}

type MemberClient interface {
	CreateMember(ctx lib.ClientContext, member types.MemberCreate) (string, error)
	List(ctx lib.ClientContext, clubID string) ([]types.Member, error)
}
type MemberClientHttp struct {
	Url string
}

func (c *MemberClientHttp) List(ctx lib.ClientContext, clubID string) ([]types.Member, error) {
	resp, err := http.Get(c.Url + "/clubs/" + clubID + "/m")
	if err != nil {
		return nil, err
	}
	var members []types.Member
	err = json.NewDecoder(resp.Body).Decode(&members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (c *MemberClientHttp) CreateMember(ctx lib.ClientContext, member types.MemberCreate) (string, error) {
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
