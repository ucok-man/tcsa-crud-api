package serializer

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONDeserializer_ErrorHandling(t *testing.T) {
	js := New()

	t.Run("malformed JSON returns error", func(t *testing.T) {
		e := echo.New()
		body := strings.NewReader(`{"name": "test"`) // Missing closing brace
		req := httptest.NewRequest(http.MethodPost, "/", body)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var result map[string]string
		err := js.Deserialize(c, &result)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "badly-formed JSON")
	})

	t.Run("wrong JSON type returns error", func(t *testing.T) {
		e := echo.New()
		body := strings.NewReader(`{"age": "not a number"}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		type User struct {
			Age int `json:"age"`
		}
		var result User
		err := js.Deserialize(c, &result)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "incorrect JSON type for field")
	})

	t.Run("empty body returns error", func(t *testing.T) {
		e := echo.New()
		body := strings.NewReader(``)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var result map[string]string
		err := js.Deserialize(c, &result)

		require.Error(t, err)
		assert.Equal(t, "body must not be empty", err.Error())
	})

	t.Run("unknown field returns error", func(t *testing.T) {
		e := echo.New()
		body := strings.NewReader(`{"name": "test", "unknown": "value"}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		type User struct {
			Name string `json:"name"`
		}
		var result User
		err := js.Deserialize(c, &result)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown key")
	})

	t.Run("multiple JSON values returns error", func(t *testing.T) {
		e := echo.New()
		body := strings.NewReader(`{"name": "test"}{"extra": "value"}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var result map[string]string
		err := js.Deserialize(c, &result)

		require.Error(t, err)
		assert.Equal(t, "body must only contain a single JSON value", err.Error())
	})

	t.Run("body too large returns error", func(t *testing.T) {
		e := echo.New()
		// Create  JSON that exceeds 1MB
		largeValue := strings.Repeat("a", 1_048_577)
		largeBody := `{"data":"` + largeValue + `"}`
		body := strings.NewReader(largeBody)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var result map[string]string
		err := js.Deserialize(c, &result)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "must not be larger than")
	})

	t.Run("valid JSON deserializes successfully", func(t *testing.T) {
		e := echo.New()
		body := strings.NewReader(`{"name": "test", "age": 30}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		type User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		var result User
		err := js.Deserialize(c, &result)

		require.NoError(t, err)
		assert.Equal(t, "test", result.Name)
		assert.Equal(t, 30, result.Age)
	})
}
