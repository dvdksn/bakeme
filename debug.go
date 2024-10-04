//go:build debug

package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func configureLogging(server *echo.Echo) {
	// Use debug mode
	server.Debug = true

	// Emit Jaeger traces
	c := jaegertracing.New(server, nil)
	defer c.Close()

	// Enable logger middleware
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:        true,
		LogURI:           true,
		LogError:         true,
		LogLatency:       true,
		LogContentLength: true,
		LogMethod:        true,
		LogRemoteIP:      true,
		LogRequestID:     true,
		LogUserAgent:     true,
		HandleError:      true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			ctx := context.Background()

			if v.Error == nil {
				// Log successful request details
				logger.LogAttrs(ctx, slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("method", v.Method),
					slog.String("remote_ip", v.RemoteIP),
					slog.Int64("latency", v.Latency.Milliseconds()),
					slog.String("content_length", v.ContentLength),
					slog.String("request_id", v.RequestID),
					slog.String("user_agent", v.UserAgent),
				)
			} else {
				// Log request details in case of an error
				logger.LogAttrs(ctx, slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("method", v.Method),
					slog.String("remote_ip", v.RemoteIP),
					slog.String("error", v.Error.Error()),
					slog.String("request_id", v.RequestID),
					slog.String("user_agent", v.UserAgent),
				)
			}
			return nil
		},
	}))
}
