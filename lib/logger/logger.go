package logger

import (
	"context"

	"go.uber.org/zap"
)

type LoggerContextKey string

const (
	key LoggerContextKey = "logger"
)

var logger *zap.Logger

func Get() *zap.Logger {
	if logger == nil {
		logger = zap.Must(zap.NewDevelopment())
	}
	return logger
}

func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(key).(*zap.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return zap.NewNop()
}

func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(key).(*zap.Logger); ok {
		if lp == l {
			return ctx
		}
	}
	return context.WithValue(ctx, key, l)
}
