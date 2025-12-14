package tlog

import (
	"bytes"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Run("shouldLog respects level setting", func(t *testing.T) {
		logger := Must(NewDevelopment())

		logger.SetLevel(log.WARN)
		assert.False(t, logger.shouldLog(log.DEBUG))
		assert.False(t, logger.shouldLog(log.INFO))
		assert.True(t, logger.shouldLog(log.WARN))
		assert.True(t, logger.shouldLog(log.ERROR))
	})

	t.Run("SetOutput changes output destination", func(t *testing.T) {
		logger := Must(NewDevelopment())
		var out bytes.Buffer

		logger.SetOutput(&out)
		logger.Print("Hello, World!")

		assert.Contains(t, out.String(), `Hello, World!`)
	})

	t.Run("withPrefix adds prefix when set", func(t *testing.T) {
		logger := Must(NewDevelopment())
		var out bytes.Buffer

		logger.SetPrefix("myapp")
		logger.SetOutput(&out)
		logger.Print("test message")

		assert.Contains(t, out.String(), `"prefix":"myapp"`)
	})

	t.Run("jsonToFields converts map to fields", func(t *testing.T) {
		logger := Must(NewDevelopment())

		j := log.JSON{
			"user_id": 123,
			"action":  "login",
		}

		fields := logger.jsonToFields(j)

		assert.Len(t, fields, 4)
		assert.Contains(t, fields, "user_id")
		assert.Contains(t, fields, 123)
		assert.Contains(t, fields, "action")
		assert.Contains(t, fields, "login")
	})

	t.Run("logJSON extracts string message correctly", func(t *testing.T) {
		logger := Must(NewDevelopment())
		logger.SetLevel(log.INFO)
		var out bytes.Buffer
		logger.SetOutput(&out)

		// With string message - should extract and use as log message
		j := log.JSON{
			"message": "user logged in",
			"user_id": 123,
		}
		logger.Printj(j)

		output := out.String()
		assert.Contains(t, output, `"msg":"user logged in"`)
		assert.Contains(t, output, `"user_id":123`)
		// Message key should be removed from fields
		assert.NotContains(t, output, `"message":"user logged in"`)
	})

	t.Run("logJSON handles missing message field", func(t *testing.T) {
		logger := Must(NewDevelopment())
		logger.SetLevel(log.INFO)
		var out bytes.Buffer
		logger.SetOutput(&out)

		// Without message - should log with empty message
		j := log.JSON{
			"user_id": 456,
			"action":  "logout",
		}
		logger.Printj(j)

		output := out.String()
		assert.Contains(t, output, `"user_id":456`)
		assert.Contains(t, output, `"action":"logout"`)
	})

	t.Run("logJSON handles non-string message", func(t *testing.T) {
		logger := Must(NewDevelopment())
		logger.SetLevel(log.INFO)
		var out bytes.Buffer
		logger.SetOutput(&out)

		// Non-string message - should keep message in fields
		j := log.JSON{
			"message": 123,
			"user_id": 789,
		}
		logger.Printj(j)

		output := out.String()
		assert.Contains(t, output, `"message":123`)
		assert.Contains(t, output, `"user_id":789`)
	})

	t.Run("shouldLog filters messages correctly", func(t *testing.T) {
		logger := Must(NewDevelopment())
		logger.SetLevel(log.WARN)
		var out bytes.Buffer
		logger.SetOutput(&out)

		// Debug should be filtered
		logger.Debug("debug message")
		assert.Empty(t, out.String())

		// Info should be filtered
		logger.Info("info message")
		assert.Empty(t, out.String())

		// Warn should pass through
		logger.Warn("warn message")
		assert.Contains(t, out.String(), "warn message")
	})
}
