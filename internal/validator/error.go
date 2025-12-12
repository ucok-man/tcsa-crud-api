package validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ucok-man/tcsa/internal/utility"
)

type ValidationErrorMap map[string]string

// Error implements the error interface and returns a formatted string of all validation errors
func (e ValidationErrorMap) Error() string {
	if len(e) == 0 {
		return "map[]"
	}

	if len(e) == 1 {
		for field, msg := range e {
			return fmt.Sprintf("%s: %s", field, msg)
		}
	}

	var builder strings.Builder
	count := 0

	for field, msg := range e {
		if count > 0 {
			builder.WriteString("; ")
		}
		builder.WriteString(fmt.Sprintf("%s: %s", field, msg))
		count++
	}

	return builder.String()
}

func (e ValidationErrorMap) MarshalJSON() ([]byte, error) {
	// Convert to regular map to use default JSON encoding
	m := map[string]string(e)

	formatted := map[string]string{}
	for key, val := range m {
		keys := strings.Split(key, ".")
		keys = utility.SlicesMap(keys, strings.ToLower)
		key = strings.Join(keys[1:], ".")

		formatted[key] = val
	}
	return json.Marshal(formatted)
}
