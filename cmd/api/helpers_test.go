package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucok-man/tcsa/internal/data"
)

func TestSortColumn(t *testing.T) {
	mockModel := new(data.MockTransactionModel)
	app := createTestApp(t, data.Models{Transactions: mockModel})

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "column with descending prefix",
			input:    "-created_at",
			expected: "created_at",
		},
		{
			name:     "column without prefix",
			input:    "created_at",
			expected: "created_at",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only dash",
			input:    "-",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.SortColumn(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSortDirection(t *testing.T) {
	mockModel := new(data.MockTransactionModel)
	app := createTestApp(t, data.Models{Transactions: mockModel})

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "descending order with dash prefix",
			input:    "-created_at",
			expected: "DESC",
		},
		{
			name:     "ascending order without prefix",
			input:    "created_at",
			expected: "ASC",
		},
		{
			name:     "empty string defaults to ascending",
			input:    "",
			expected: "ASC",
		},
		{
			name:     "only dash returns descending",
			input:    "-",
			expected: "DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.SortDirection(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPageOffset(t *testing.T) {
	mockModel := new(data.MockTransactionModel)
	app := createTestApp(t, data.Models{Transactions: mockModel})

	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{
			name:     "first page with page size 10",
			page:     1,
			pageSize: 10,
			expected: 0,
		},
		{
			name:     "second page with page size 10",
			page:     2,
			pageSize: 10,
			expected: 10,
		},
		{
			name:     "third page with page size 20",
			page:     3,
			pageSize: 20,
			expected: 40,
		},
		{
			name:     "page 5 with page size 25",
			page:     5,
			pageSize: 25,
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.PageOffset(tt.page, tt.pageSize)
			assert.Equal(t, tt.expected, result)
		})
	}
}
