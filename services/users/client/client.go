package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type UserClient struct {
	Url string
}

func (c UserClient) GetUsers() ([]User, error) {
	resp, err := http.Get(c.Url)
	if err != nil {
		return nil, err
	}
	var members []User
	err = json.NewDecoder(resp.Body).Decode(&members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (c UserClient) CreateUser(member UserCreate) (string, error) {
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
