package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ucok-man/tcsa/internal/tlog"
	"go.uber.org/zap"
)

func (app *application) withRecover() echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			c.Logger().Error("Recovering from panic",
				zap.Error(err),
				zap.Any("url", c.Request().URL),
				zap.String("method", c.Request().Method),
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
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			switch {
			case v.Status >= 500:
				c.Logger().Errorj(tlog.JSON{
					"message":       http.StatusText(v.Status),
					"code":          v.Status,
					"method":        v.Method,
					"url":           v.URI,
					"ip_addr":       v.RemoteIP,
					"response_time": v.Latency,
					"response_size": v.ResponseSize,
					"error":         v.Error,
				})
			default:
				c.Logger().Infoj(tlog.JSON{
					"message":       http.StatusText(v.Status),
					"code":          v.Status,
					"method":        v.Method,
					"url":           v.URI,
					"ip_addr":       v.RemoteIP,
					"response_time": v.Latency,
					"response_size": v.ResponseSize,
				})

			}

			return nil
		},
	})
}
