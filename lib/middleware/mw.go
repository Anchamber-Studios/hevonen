package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	ory "github.com/ory/client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func AuthPaseto(tokenKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "missing authorization header")
			}
			// _ := auth[len("Bearer "):]
			return next(c)
		}
	}
}

func AuthJWTOry() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenHeader := c.Request().Header.Get("Authorization")
			if tokenHeader == "" {
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "missing authorization header")
			}

			config := ory.NewConfiguration()
			config.Servers = ory.ServerConfigurations{{URL: fmt.Sprintf("http://localhost:%s/.ory", "4000")}}

			client := ory.NewAPIClient(config)
			req := client.OAuth2API.IntrospectOAuth2Token(c.Request().Context()).
				Token(tokenHeader[len("Bearer "):])
			token, _, err := client.OAuth2API.IntrospectOAuth2TokenExecute(req)
			if err != nil {
				c.Logger().Infof("error introspecting token: %s", err)
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "invalid token")
			}
			if !token.Active {
				c.Logger().Warnf("token is not active")
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "invalid token")
			}
			return next(c)
		}
	}
}

func Logging(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			fields := []zapcore.Field{
				zap.String("remote_ip", c.RealIP()),
				zap.String("latency", time.Since(start).String()),
				zap.String("host", req.Host),
				zap.String("request", fmt.Sprintf("%s %s", req.Method, req.RequestURI)),
				zap.Int("status", res.Status),
				zap.Int64("size", res.Size),
				zap.String("user_agent", req.UserAgent()),
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			fields = append(fields, zap.String("request_id", id))

			// n := res.Status
			// switch {
			// case n >= 500:
			// 	log.With(zap.Error(err)).Error("Server error", fields...)
			// case n >= 400:
			// 	log.With(zap.Error(err)).Warn("Client error", fields...)
			// case n >= 300:
			// 	log.Info("Redirection", fields...)
			// default:
			// 	log.Info("Success", fields...)
			// }

			return nil
		}
	}
}
