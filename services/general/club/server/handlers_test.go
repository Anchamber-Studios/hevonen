package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/services/club/client"
	repo "github.com/anchamber-studios/hevonen/services/club/db"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func TestCreateUser(t *testing.T) {
	e := echo.New()
	m := testMemberCreate()
	memberRepo := &repo.MemberRepoMock{}
	data, _ := json.Marshal(m)
	req := httptest.NewRequest(http.MethodPost, "/members", bytes.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	hander := MemberHandler{}
	if err := hander.new(ctx(t, c, Repos{Members: memberRepo})); err != nil {
		t.Errorf("Unable to create member: %v\n", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/members", bytes.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if err := hander.new(ctx(t, c, Repos{Members: memberRepo})); err == nil {
		t.Errorf("should return error for creating user with same mail twice: %v\n", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/members", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if err := hander.list(ctx(t, c, Repos{Members: memberRepo})); err != nil {
		t.Errorf("Unable to list members: %v\n", err)
	}
	var members []client.Member
	json.NewDecoder(rec.Body).Decode(&members)
	if len(members) != 1 {
		t.Errorf("Expected 1 member, got %d\n", len(members))
	}
}

func testConfig() config.Config {
	return config.Config{
		Port: "8080",
		Host: "localhost",
		Tls: config.TlsConfig{
			Enabled: false,
			Key:     "",
			Cert:    "",
		},
		Auth: config.Auth{
			ClientId:     "test",
			ClientSecret: "test",
		},
		DB: config.DB{
			Url:      "localhost",
			Port:     "5432",
			Database: "members",
			User:     "members",
			Password: "members",
		},
	}
}

func testDb() *pgx.Conn {
	return nil
}

func ctx(t *testing.T, c echo.Context, repos Repos) *CustomContext {
	idConv, err := setupIdConversion()
	if err != nil {
		t.Errorf("Unable to setup id conversion: %v\n", err)
	}
	return &CustomContext{c, testConfig(), testDb(), idConv, repos, Services{}}
}

func testMemberCreate() client.MemberCreate {
	return client.MemberCreate{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jd@hevonen.io",
		Phone:     "1234567890",
	}
}
