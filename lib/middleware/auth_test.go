package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/anchamber-studios/hevonen/lib/middleware"
	"github.com/labstack/echo/v4"
)

func TestAuthPaseto(t *testing.T) {
	k := "aaabbbccc0123456789012345678901212345678901234567890123456789012"
	key, err := paseto.V4SymmetricKeyFromHex(k)
	if err != nil {
		t.Errorf("failed to create test token")
	}
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))
	signed := token.V4Encrypt(key, nil)
	testCases := []struct {
		name    string
		token   string
		success bool
	}{
		{
			name:    "With signed token",
			token:   signed,
			success: true,
		},
		{
			name:    "With empty token",
			token:   "",
			success: false,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mw := middleware.AuthPaseto(k)
			handler := mw(func(c echo.Context) error {
				return nil
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", tc.token))
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			err := handler(c)

			if err != nil && tc.success {
				t.Errorf("%s: AuthPaseto() = %v, want nil", tc.name, err)
			} else if err == nil && !tc.success {
				t.Errorf("%s: AuthPaseto() = nil, want error", tc.name)
			}
		})
	}
}
