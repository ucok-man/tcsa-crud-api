package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlicesMapInt(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		mapFunc  mapFunc[int]
		expected []int
	}{
		{
			name:     "double all elements",
			input:    []int{1, 2, 3, 4},
			mapFunc:  func(x int) int { return x * 2 },
			expected: []int{2, 4, 6, 8},
		},
		{
			name:     "empty slice",
			input:    []int{},
			mapFunc:  func(x int) int { return x * 2 },
			expected: []int{},
		},
		{
			name:     "negate all elements",
			input:    []int{1, -2, 3, -4},
			mapFunc:  func(x int) int { return -x },
			expected: []int{-1, 2, -3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SlicesMap(tt.input, tt.mapFunc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSlicesMapString(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		mapFunc  mapFunc[string]
		expected []string
	}{
		{
			name:     "uppercase strings",
			input:    []string{"hello", "world"},
			mapFunc:  func(s string) string { return s + "!" },
			expected: []string{"hello!", "world!"},
		},
		{
			name:     "empty string slice",
			input:    []string{},
			mapFunc:  func(s string) string { return s + "!" },
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SlicesMap(tt.input, tt.mapFunc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSetPtrValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "integer value",
			input:    42,
			expected: 42,
		},
		{
			name:     "string value",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "float value",
			input:    3.14,
			expected: 3.14,
		},
		{
			name:     "zero value",
			input:    0,
			expected: 0,
		},
		{
			name:     "boolean value",
			input:    true,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				result := SetPtrValue(v)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected, *result)
			case string:
				result := SetPtrValue(v)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected, *result)
			case float64:
				result := SetPtrValue(v)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected, *result)
			case bool:
				result := SetPtrValue(v)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected, *result)
			}
		})
	}
}

func TestDerefOrDefault(t *testing.T) {
	intVal := 42
	strVal := "hello"
	floatVal := 3.14

	tests := []struct {
		name       string
		ptr        any
		defaultVal any
		expected   any
	}{
		{
			name:       "non-nil int pointer",
			ptr:        &intVal,
			defaultVal: 0,
			expected:   42,
		},
		{
			name:       "nil int pointer",
			ptr:        (*int)(nil),
			defaultVal: 99,
			expected:   99,
		},
		{
			name:       "non-nil string pointer",
			ptr:        &strVal,
			defaultVal: "default",
			expected:   "hello",
		},
		{
			name:       "nil string pointer",
			ptr:        (*string)(nil),
			defaultVal: "default",
			expected:   "default",
		},
		{
			name:       "non-nil float pointer",
			ptr:        &floatVal,
			defaultVal: 0.0,
			expected:   3.14,
		},
		{
			name:       "nil float pointer",
			ptr:        (*float64)(nil),
			defaultVal: 1.5,
			expected:   1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch ptr := tt.ptr.(type) {
			case *int:
				result := DerefOrDefault(ptr, tt.defaultVal.(int))
				assert.Equal(t, tt.expected, result)
			case *string:
				result := DerefOrDefault(ptr, tt.defaultVal.(string))
				assert.Equal(t, tt.expected, result)
			case *float64:
				result := DerefOrDefault(ptr, tt.defaultVal.(float64))
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestRound2(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "round up",
			input:    3.456,
			expected: 3.46,
		},
		{
			name:     "round down",
			input:    3.454,
			expected: 3.45,
		},
		{
			name:     "already rounded",
			input:    3.45,
			expected: 3.45,
		},
		{
			name:     "zero",
			input:    0.0,
			expected: 0.0,
		},
		{
			name:     "negative round up",
			input:    -3.456,
			expected: -3.46,
		},
		{
			name:     "negative round down",
			input:    -3.454,
			expected: -3.45,
		},
		{
			name:     "very small number",
			input:    0.00001,
			expected: 0.00,
		},
		{
			name:     "round 0.005",
			input:    0.005,
			expected: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Round2(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
