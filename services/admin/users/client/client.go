package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/anchamber-studios/hevonen/lib"
)

type UserClient interface {
	GetUsers(ctx lib.ClientContext) ([]User, error)
	GetUser(ctx lib.ClientContext, id string) (User, error)
	Register(ctx lib.ClientContext, user UserCreate) (string, error)
	Login(ctx lib.ClientContext, login UserLogin) (UserLoginResponse, error)
	Logout(ctx lib.ClientContext, id string) error
}

type UserClientHttp struct {
	Url string
}

func (c *UserClientHttp) GetUsers(ctx lib.ClientContext) ([]User, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.Url, nil)
	if err != nil {
		return nil, err
	}
	ctx.SetHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var contacts []User
	err = json.NewDecoder(resp.Body).Decode(&contacts)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (c *UserClientHttp) GetUser(ctx lib.ClientContext, id string) (User, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.Url, id), nil)
	if err != nil {
		return User{}, err
	}
	ctx.SetHeader(req)
	resp, err := client.Do(req)
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

func (c *UserClientHttp) Register(ctx lib.ClientContext, user UserCreate) (string, error) {
	client := &http.Client{}

	userJson, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", c.Url+"/register", bytes.NewReader(userJson))
	if err != nil {
		return "", err
	}
	ctx.SetHeader(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	location := res.Header.Get("Location")
	return location, nil
}

func (c *UserClientHttp) Login(ctx lib.ClientContext, login UserLogin) (UserLoginResponse, error) {
	client := &http.Client{}

	loginJson, err := json.Marshal(login)
	if err != nil {
		return UserLoginResponse{}, err
	}
	req, err := http.NewRequest("POST", c.Url+"/login", bytes.NewReader(loginJson))
	if err != nil {
		log.Printf("req err: %v\n", err)
		return UserLoginResponse{}, err
	}
	ctx.SetHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("resp err: %v\n", err)
		return UserLoginResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return UserLoginResponse{}, lib.ErrUnauthorized
	}
	var loginResponse UserLoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	if err != nil {
		return UserLoginResponse{}, err
	}
	return loginResponse, nil
}

func (c *UserClientHttp) Logout(ctx lib.ClientContext, id string) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/logout", c.Url, id), nil)
	if err != nil {
		return err
	}
	ctx.SetHeader(req)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
