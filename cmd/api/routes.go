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
	ec.HTTPErrorHandler = app.HTTPErrorHandler

	ec.Use(app.withRecover())
	ec.Use(app.withRequestLogger())

	// Documentation routes

	ec.FileFS("/swagger.yaml", "docs/swagger.yaml", swaggerFile)
	ec.GET("/docs", app.serveSwaggerUI)

	// Health check
	ec.GET("/healthcheck", app.healthcheckHandler)

	// Transaction routes
	transactions := ec.Group("/transactions")
	{
		transactions.GET("", app.getAllTransactionHandler)
		transactions.POST("", app.createTransactionHandler)
		transactions.GET("/:id", app.getByIdTransactionHandler)
		transactions.PUT("/:id", app.updateByIdTransactionHandler)
		transactions.DELETE("/:id", app.removeByIdTransactionHandler)
	}

	// Dashboard routes
	dashboard := ec.Group("/dashboard")
	{
		dashboard.GET("/summary", app.summaryTransactionHandler)
	}

	return ec
}
