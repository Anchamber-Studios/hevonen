package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anchamber-studios/hevonen/lib"
	"github.com/anchamber-studios/hevonen/lib/config"
	ldb "github.com/anchamber-studios/hevonen/lib/db"
	"github.com/anchamber-studios/hevonen/services/club/db"
	"github.com/anchamber-studios/hevonen/services/club/server"
	"github.com/anchamber-studios/hevonen/services/club/services"
	"github.com/anchamber-studios/hevonen/services/club/shared/types"
	"github.com/labstack/echo/v4"
)

const (
	TestIdentityID = "b4df50de-7eca-41bd-85b1-08acfe61403d"
	TestEmail      = "test@hevonen.io"
)

func createContext(t *testing.T, r *http.Request, w http.ResponseWriter, config config.Config) server.CustomContext {
	e := echo.New()
	err := db.SetupDB(config)
	if err != nil {
		t.Fatalf("Unable to setup database: %v\n", err)
	}
	conn, err := ldb.OpenConnection(config)
	if err != nil {
		t.Fatalf("Unable to connect to database: %v\n", err)
	}
	idc, err := server.SetupIdConversion()
	if err != nil {
		t.Fatalf("Unable setup id conversion: %v\n", err)
	}

	ctx := server.CustomContext{
		e.NewContext(r, w),
		config,
		server.Services{
			Clubs: services.NewClubService(&db.ClubRepoPostgre{DB: conn, IdConversion: idc}),
		},
	}
	ctx.Set("identityID", TestIdentityID)
	ctx.Set("email", TestEmail)

	return ctx
}

func creaetConfig() config.Config {
	config := config.DefaultConfig()
	// config.DB.Database = "club_" + time.Now().Format("_2006-01-02_15-04-05")
	return config
}

func TestClubHandlerCreate(t *testing.T) {
	conf := creaetConfig()

	defer ldb.ResetDB(context.Background(), "clubs", conf)
	handler := &server.ClubHandler{}
	rec, err := createClub(t, handler)
	if err != nil {
		t.Fatalf("Error creating club: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rec.Code)
	}
	if rec.Header().Get("Location") == "" {
		t.Errorf("Expected Location header to be set")
	}

	id := rec.Header().Get("Location")[3:]
	deleteClub(t, handler, conf, id)
}

func TestClubHandlerCreateClubWithAlreadyUsedName(t *testing.T) {
	conf := creaetConfig()

	defer ldb.ResetDB(context.Background(), "clubs", conf)
	handler := &server.ClubHandler{}
	rec, err := createClub(t, handler)
	if err != nil {
		t.Fatalf("Error creating club: %v", err)
	}
	id := rec.Header().Get("Location")[3:]
	defer deleteClub(t, handler, conf, id)

	_, err = createClub(t, handler)
	if err == nil {
		t.Fatalf("Should return error when name is already used: %v", err)
	}
	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("Expected an echo.HTTPError, got %T", err)
	}
	if httpErr.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", httpErr.Code)
	}
	apiErr, ok := httpErr.Message.(*lib.ApiError)
	if !ok {
		t.Fatalf("Expected an lib.ApiError, got %T", httpErr.Message)
	}
	if apiErr.StatusCode != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", apiErr.StatusCode)
	}

}

func deleteClub(t *testing.T, handler *server.ClubHandler, conf config.Config, id string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Add("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := createContext(t, req, rec, conf)
	c.SetPath("/c/:clubID")
	c.SetParamNames("clubID")
	c.SetParamValues(id)
	if err := handler.DeleteClub(&c); err != nil {
		return rec, err
	}
	return rec, nil
}

func createClub(t *testing.T, handler *server.ClubHandler) (*httptest.ResponseRecorder, error) {
	club := types.ClubCreate{Name: "Test Club", Email: "test@hevonen.io"}
	data, err := json.Marshal(club)
	if err != nil {
		t.Fatalf("Error marshalling club: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	conf := creaetConfig()
	c := createContext(t, req, rec, conf)
	if err := handler.Create(&c); err != nil {
		return rec, err
	}
	return rec, nil
}
