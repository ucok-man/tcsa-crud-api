package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ucok-man/tcsa/internal/data"
	"github.com/ucok-man/tcsa/internal/serializer"
	"github.com/ucok-man/tcsa/internal/tlog"
	"github.com/ucok-man/tcsa/internal/validator"
)

// createTestApp creates a test application with mocked dependencies
func createTestApp(t *testing.T, mock data.Models) *application {
	t.Helper()

	logger := tlog.Must(tlog.NewDevelopment())
	logger.SetOutput(&bytes.Buffer{})

	return &application{
		config: Config{
			Port: 3000,
			Env:  "test",
		},
		logger: logger,
		models: mock,
	}
}

// createTestContext creates a new Echo context for testing
func createTestContext(method, path string, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.JSONSerializer = serializer.New()
	e.Validator = validator.New()

	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}
