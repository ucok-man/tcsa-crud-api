package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucok-man/tcsa/internal/serializer"
	"github.com/ucok-man/tcsa/internal/validator"
)

func (app *application) routes() http.Handler {
	ec := echo.New()
	ec.JSONSerializer = serializer.New()
	ec.Validator = validator.New()
	ec.Logger = app.logger
	// ec.HTTPErrorHandler

	ec.Use(app.withRecover())
	ec.Use(app.withRequestLogger())

	ec.GET("/healthcheck", app.healthcheckHandler)

	return ec
}
