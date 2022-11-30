package logger

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func InitializeLogger() *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	root := zerolog.New(os.Stdout)
	root = root.With().Caller().Logger()
	return &root
}

type ctxKey struct{}

var loggerKey = &ctxKey{}

func LoggerMW(ctx context.Context, logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, loggerKey, logger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func From(ctx context.Context) *zerolog.Logger {
	service := ctx.Value(loggerKey)
	logger, ok := service.(*zerolog.Logger)
	if logger == nil || !ok {
		return nil
	}
	return logger
}
