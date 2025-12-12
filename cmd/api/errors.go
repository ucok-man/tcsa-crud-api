package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucok-man/tcsa/internal/tlog"
)

func (app *application) ErrInternalServer(err error, message string, req *http.Request) error {
	app.logger.Errorj(tlog.JSON{
		"message": message,
		"path":    req.URL,
		"method":  req.Method,
		"error":   err,
	})
	return echo.NewHTTPError(
		http.StatusInternalServerError,
		"the server encountered a problem and could not process your request",
	)
}

func (app *application) ErrNotFound(customMsg ...string) error {
	message := "the requested resource could not be found"
	if len(customMsg) > 0 && customMsg[0] != "" {
		message = customMsg[0]
	}
	return echo.NewHTTPError(http.StatusNotFound, message)
}

func (app *application) ErrMethodNotAllowed(method string) error {
	return echo.NewHTTPError(
		http.StatusMethodNotAllowed,
		fmt.Sprintf("the %s method is not supported for this resource", method),
	)
}

func (app *application) ErrBadRequest(message string) error {
	return echo.NewHTTPError(http.StatusBadRequest, message)
}

func (app *application) ErrFailedValidation(errmap any) error {
	return echo.NewHTTPError(http.StatusUnprocessableEntity, errmap)
}

func (app *application) ErrEditConflict() error {
	return echo.NewHTTPError(
		http.StatusConflict,
		"unable to update the record due to an edit conflict, please try again",
	)
}

func (app *application) ErrRateLimitExceeded() error {
	return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
}

func (app *application) ErrForbidden(message ...string) error {
	msg := "forbidden"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return echo.NewHTTPError(http.StatusForbidden, msg)
}
