package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/anchamber-studios/hevonen/lib"
)

type ProfileClient interface {
	CreateProfile(ctx lib.ClientContext, profile ProfileCreateRequest) (string, error)
	GetProfile(ctx lib.ClientContext, profileId string) (ProfileResponse, error)
	UpdateProfile(ctx lib.ClientContext, profileId string, profile ProfileUpdateRequest) error
	// DeleteProfile(ctx lib.ClientContext, profileId string) error
	// AddAddress(ctx lib.ClientContext, profileId string, address AddressCreateRequest) (string, error)
	// GetAddresses(ctx lib.ClientContext, profileId string) ([]AddressResponse, error)
	// DeleteAddress(ctx lib.ClientContext, profileId string, addressId string) error
	// AddContactInfo(ctx lib.ClientContext, profileId string, contactInfo ContactInformationCreateRequest) (string, error)
	// GetContactInfo(ctx lib.ClientContext, profileId string) ([]ContactInformationResponse, error)
	// DeleteContactInfo(ctx lib.ClientContext, profileId string, contactInfoId string) error
}

type ProfileClientHttp struct {
	Url string
}

func (c *ProfileClientHttp) GetProfile(ctx lib.ClientContext, profileId string) (ProfileResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.Url, profileId), nil)
	if err != nil {
		return ProfileResponse{}, err
	}
	ctx.SetHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return ProfileResponse{}, err
	}
	if resp.StatusCode/100 != 2 {
		apiErr := &lib.ApiError{}
		json.NewDecoder(resp.Body).Decode(apiErr)
		fmt.Printf("err: %v\n", apiErr)
		return ProfileResponse{}, apiErr
	}
	var profile ProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return ProfileResponse{}, err
	}
	return profile, nil
}

func (c *ProfileClientHttp) CreateProfile(ctx lib.ClientContext, profile ProfileCreateRequest) (string, error) {
	client := &http.Client{}

	profileJson, err := json.Marshal(profile)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", c.Url, bytes.NewReader(profileJson))
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
	locationComponents := strings.Split(location, "/")
	return locationComponents[len(locationComponents)-1], nil
}

func (c *ProfileClientHttp) UpdateProfile(ctx lib.ClientContext, profileId string, profile ProfileUpdateRequest) error {
	client := &http.Client{}

	profileJson, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	fmt.Printf("profileJson: %s\n", profileJson)
	fmt.Printf("profile: %v\n", profile)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s", c.Url, profileId), bytes.NewReader(profileJson))
	if err != nil {
		return err
	}
	ctx.SetHeader(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode/100 != 2 {
		apiErr := &lib.ApiError{}
		json.NewDecoder(res.Body).Decode(apiErr)
		return apiErr
	}
	return nil
}
