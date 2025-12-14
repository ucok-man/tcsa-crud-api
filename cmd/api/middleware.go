package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ucok-man/tcsa/internal/tlog"
	"go.uber.org/zap"
)

func (app *application) withRecover() echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(ctx echo.Context, err error, stack []byte) error {
			ctx.Logger().Error("Recovering from panic",
				zap.Error(err),
				zap.Any("url", ctx.Request().URL),
				zap.String("method", ctx.Request().Method),
			)
			return err
		},
	})
}

func (app *application) withRequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRemoteIP:     true,
		LogStatus:       true,
		LogMethod:       true,
		LogURI:          true,
		LogLatency:      true,
		LogResponseSize: true,
		LogError:        true,
		LogValuesFunc: func(ctx echo.Context, v middleware.RequestLoggerValues) error {
			var message string
			switch {
			case v.Status >= 300 && v.Status <= 399:
				message = "Redirect"
			case v.Status >= 400 && v.Status <= 499:
				message = "Client Error"
			case v.Status >= 500 || v.Error != nil:
				v.Status = 500
				message = "Server Error"
			default:
				message = "Success"
			}

			data := tlog.JSON{
				"message":       message,
				"code":          v.Status,
				"method":        v.Method,
				"url":           v.URI,
				"ip_addr":       v.RemoteIP,
				"response_time": v.Latency,
				"response_size": v.ResponseSize,
			}

			ctx.Logger().Infoj(data)

			return nil
		},
	})
}

func (app *application) withCORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: app.config.Cors.TrustedOrigins,
	})
}
