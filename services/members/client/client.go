package client

import (
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
