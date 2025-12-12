package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) healthcheckHandler(ctx echo.Context) error {
	env := envelope{
		"status": "available",
		"system_info": map[string]any{
			"environment": app.config.Env,
			"version":     VERSION,
		},
	}

	return ctx.JSON(http.StatusOK, &env)
}
