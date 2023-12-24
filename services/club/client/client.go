package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type MemberClient struct {
	Url string
}

func (c MemberClient) GetMembers() ([]Member, error) {
	resp, err := http.Get(c.Url)
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

func (c MemberClient) CreateMember(member MemberCreate) (string, error) {
	memberJson, err := json.Marshal(member)
	if err != nil {
		return "", err
	}
	res, err := http.Post(c.Url, "application/json", bytes.NewReader(memberJson))
	if err != nil {
		return "", err
	}
	location := res.Header.Get("Location")
	return location, nil
}
