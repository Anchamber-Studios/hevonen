package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/anchamber-studios/hevonen/lib"
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

func (c UserClient) CreateUser(user UserCreate) (string, error) {
	userJson, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	res, err := http.Post(c.Url, "application/json", bytes.NewReader(userJson))
	if err != nil {
		return "", err
	}
	location := res.Header.Get("Location")
	return location, nil
}

func (c UserClient) GetUser(id string) (User, error) {
	resp, err := http.Get(c.Url + "/" + id)
	if err != nil {
		return User{}, err
	}
	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (c UserClient) Login(login UserLogin) (User, error) {
	loginJson, err := json.Marshal(login)
	if err != nil {
		return User{}, err
	}
	resp, err := http.Post(c.Url+"/login", "application/json", bytes.NewReader(loginJson))
	if err != nil {
		return User{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return User{}, lib.ErrUnauthorized
	}
	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
