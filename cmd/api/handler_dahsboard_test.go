package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ucok-man/tcsa/internal/data"
)

func TestSummaryTransactionHandler(t *testing.T) {
	t.Run("successfully gets transaction summary with default pagination", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedTransactions := []*data.Transaction{
			{
				ID:        1,
				UserId:    1,
				Amount:    10000,
				Status:    data.TransactionStatusSucces,
				Version:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        2,
				UserId:    1,
				Amount:    5000,
				Status:    data.TransactionStatusPending,
				Version:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        3,
				UserId:    2,
				Amount:    8000,
				Status:    data.TransactionStatusFailed,
				Version:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		expectedSummary := &data.TransactionSummary{
			Transactions: expectedTransactions,
			Summary: data.Summary{
				CountTotal: 3,
				Pending: struct {
					Count          int     `json:"count"`
					RatePercentage float64 `json:"rate_percentage"`
				}{
					Count:          1,
					RatePercentage: 33.33,
				},
				Success: struct {
					Count          int     `json:"count"`
					RatePercentage float64 `json:"rate_percentage"`
				}{
					Count:          1,
					RatePercentage: 33.33,
				},
				Failed: struct {
					Count          int     `json:"count"`
					RatePercentage float64 `json:"rate_percentage"`
				}{
					Count:          1,
					RatePercentage: 33.34,
				},
			},
		}

		expectedMetadata := &data.Metadata{
			CurrentPage:  1,
			PageSize:     10,
			FirstPage:    1,
			LastPage:     1,
			TotalRecords: 3,
		}

		mockModel.On("Summary", mock.MatchedBy(func(param data.TransactionSummaryParam) bool {
			return param.Page == 1 &&
				param.PageSize == 10 &&
				param.SortColumn == "id" &&
				param.SortDirection == "ASC" &&
				param.FilterDateRange == 0 &&
				param.FilterUserId == 0
		})).Return(expectedSummary, expectedMetadata, nil)

		ctx, rec := createTestContext(http.MethodGet, "/dashboard/summary", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		// Check data exists
		summaryData := response["data"].(map[string]interface{})
		assert.NotNil(t, summaryData)

		// Check transactions
		transactions := summaryData["transactions"].([]interface{})
		assert.Len(t, transactions, 3)

		// Check summary stats
		summary := summaryData["summary"].(map[string]interface{})
		assert.Equal(t, float64(3), summary["total_transaction"])

		pending := summary["pending"].(map[string]interface{})
		assert.Equal(t, float64(1), pending["count"])
		assert.Equal(t, 33.33, pending["rate_percentage"])

		success := summary["success"].(map[string]interface{})
		assert.Equal(t, float64(1), success["count"])
		assert.Equal(t, 33.33, success["rate_percentage"])

		failed := summary["failed"].(map[string]interface{})
		assert.Equal(t, float64(1), failed["count"])
		assert.Equal(t, 33.34, failed["rate_percentage"])

		// Check metadata
		metadata := response["metadata"].(map[string]interface{})
		assert.Equal(t, float64(1), metadata["current_page"])
		assert.Equal(t, float64(10), metadata["page_size"])
		assert.Equal(t, float64(3), metadata["total_records"])

		mockModel.AssertExpectations(t)
	})

	t.Run("successfully gets summary with custom pagination", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedSummary := &data.TransactionSummary{
			Transactions: []*data.Transaction{},
			Summary: data.Summary{
				CountTotal: 50,
				Pending: struct {
					Count          int     `json:"count"`
					RatePercentage float64 `json:"rate_percentage"`
				}{Count: 10, RatePercentage: 20.0},
				Success: struct {
					Count          int     `json:"count"`
					RatePercentage float64 `json:"rate_percentage"`
				}{Count: 30, RatePercentage: 60.0},
				Failed: struct {
					Count          int     `json:"count"`
					RatePercentage float64 `json:"rate_percentage"`
				}{Count: 10, RatePercentage: 20.0},
			},
		}

		expectedMetadata := &data.Metadata{
			CurrentPage:  2,
			PageSize:     20,
			FirstPage:    1,
			LastPage:     3,
			TotalRecords: 50,
		}

		mockModel.On("Summary", mock.MatchedBy(func(param data.TransactionSummaryParam) bool {
			return param.Page == 2 && param.PageSize == 20
		})).Return(expectedSummary, expectedMetadata, nil)

		ctx, rec := createTestContext(http.MethodGet, "/dashboard/summary?page=2&page_size=20", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		metadata := response["metadata"].(map[string]interface{})
		assert.Equal(t, float64(2), metadata["current_page"])
		assert.Equal(t, float64(20), metadata["page_size"])
		assert.Equal(t, float64(50), metadata["total_records"])

		mockModel.AssertExpectations(t)
	})

	t.Run("successfully filters by date range", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedSummary := &data.TransactionSummary{
			Transactions: []*data.Transaction{},
			Summary: data.Summary{
				CountTotal: 10,
				Pending: struct {
					Count          int     `json:"count"`
					RatePercentage float64 `json:"rate_percentage"`
				}{Count: 3, RatePercentage: 30.0},
				Success: struct {
					Count          int     `json:"count"`
					RatePercentage float64 `json:"rate_percentage"`
				}{Count: 5, RatePercentage: 50.0},
				Failed: struct {
					Count          int     `json:"count"`
					RatePercentage float64 `json:"rate_percentage"`
				}{Count: 2, RatePercentage: 20.0},
			},
		}

		expectedMetadata := &data.Metadata{
			CurrentPage:  1,
			PageSize:     10,
			FirstPage:    1,
			LastPage:     1,
			TotalRecords: 10,
		}

		mockModel.On("Summary", mock.MatchedBy(func(param data.TransactionSummaryParam) bool {
			return param.FilterDateRange == 30 // Last 30 days
		})).Return(expectedSummary, expectedMetadata, nil)

		ctx, rec := createTestContext(http.MethodGet, "/dashboard/summary?date_range=30", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		summaryData := response["data"].(map[string]interface{})
		summary := summaryData["summary"].(map[string]interface{})
		assert.Equal(t, float64(10), summary["total_transaction"])

		mockModel.AssertExpectations(t)
	})

	t.Run("successfully filters by user id", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedSummary := &data.TransactionSummary{
			Transactions: []*data.Transaction{
				{
					ID:     1,
					UserId: 5,
					Amount: 10000,
					Status: data.TransactionStatusSucces,
				},
				{
					ID:     2,
					UserId: 5,
					Amount: 5000,
					Status: data.TransactionStatusSucces,
				},
			},
			Summary: data.Summary{
				CountTotal: 2,
				Success: struct {
					Count          int     `json:"count"`
					RatePercentage float64 `json:"rate_percentage"`
				}{Count: 2, RatePercentage: 100.0},
			},
		}

		expectedMetadata := &data.Metadata{
			CurrentPage:  1,
			PageSize:     10,
			FirstPage:    1,
			LastPage:     1,
			TotalRecords: 2,
		}

		mockModel.On("Summary", mock.MatchedBy(func(param data.TransactionSummaryParam) bool {
			return param.FilterUserId == 5
		})).Return(expectedSummary, expectedMetadata, nil)

		ctx, rec := createTestContext(http.MethodGet, "/dashboard/summary?user_id=5", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		summaryData := response["data"].(map[string]interface{})
		transactions := summaryData["transactions"].([]interface{})
		assert.Len(t, transactions, 2)

		mockModel.AssertExpectations(t)
	})

	t.Run("successfully applies custom sorting", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedSummary := &data.TransactionSummary{
			Transactions: []*data.Transaction{},
			Summary:      data.Summary{CountTotal: 5},
		}

		expectedMetadata := &data.Metadata{
			CurrentPage:  1,
			PageSize:     10,
			FirstPage:    1,
			LastPage:     1,
			TotalRecords: 5,
		}

		// Test descending sort by amount
		mockModel.On("Summary", mock.MatchedBy(func(param data.TransactionSummaryParam) bool {
			return param.SortColumn == "amount" && param.SortDirection == "DESC"
		})).Return(expectedSummary, expectedMetadata, nil)

		ctx, rec := createTestContext(http.MethodGet, "/dashboard/summary?sort_by=-amount", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockModel.AssertExpectations(t)
	})

	t.Run("successfully combines multiple filters", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedSummary := &data.TransactionSummary{
			Transactions: []*data.Transaction{},
			Summary: data.Summary{
				CountTotal: 3,
				Success: struct {
					Count          int     `json:"count"`
					RatePercentage float64 `json:"rate_percentage"`
				}{Count: 3, RatePercentage: 100.0},
			},
		}

		expectedMetadata := &data.Metadata{
			CurrentPage:  1,
			PageSize:     10,
			FirstPage:    1,
			LastPage:     1,
			TotalRecords: 3,
		}

		mockModel.On("Summary", mock.MatchedBy(func(param data.TransactionSummaryParam) bool {
			return param.FilterDateRange == 7 &&
				param.FilterUserId == 3 &&
				param.SortColumn == "created_at" &&
				param.SortDirection == "DESC"
		})).Return(expectedSummary, expectedMetadata, nil)

		ctx, rec := createTestContext(
			http.MethodGet,
			"/dashboard/summary?date_range=7&user_id=3&sort_by=-created_at",
			"",
		)

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockModel.AssertExpectations(t)
	})

	t.Run("returns validation error for invalid page number", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/dashboard/summary?page=0", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Summary")
	})

	t.Run("returns validation error for page number exceeding maximum", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/dashboard/summary?page=1001", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Summary")
	})

	t.Run("returns validation error for invalid page size", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/dashboard/summary?page_size=0", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Summary")
	})

	t.Run("returns validation error for page size exceeding maximum", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/dashboard/summary?page_size=101", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Summary")
	})

	t.Run("returns validation error for invalid sort_by value", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/dashboard/summary?sort_by=invalid_field", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Summary")
	})

	t.Run("returns validation error for invalid date_range", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/dashboard/summary?date_range=0", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Summary")
	})

	t.Run("returns validation error for date_range exceeding maximum", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/dashboard/summary?date_range=367", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Summary")
	})

	t.Run("returns validation error for invalid user_id", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/dashboard/summary?user_id=0", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Summary")
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		mockModel.On("Summary", mock.AnythingOfType("data.TransactionSummaryParam")).
			Return(nil, nil, assert.AnError)

		ctx, _ := createTestContext(http.MethodGet, "/dashboard/summary", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertExpectations(t)
	})

	t.Run("handles empty result set gracefully", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedSummary := &data.TransactionSummary{
			Transactions: []*data.Transaction{},
			Summary: data.Summary{
				CountTotal: 0,
			},
		}

		expectedMetadata := &data.Metadata{}

		mockModel.On("Summary", mock.AnythingOfType("data.TransactionSummaryParam")).
			Return(expectedSummary, expectedMetadata, nil)

		ctx, rec := createTestContext(http.MethodGet, "/dashboard/summary", "")

		// Execute
		err := app.summaryTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		summaryData := response["data"].(map[string]interface{})
		transactions := summaryData["transactions"].([]interface{})
		assert.Empty(t, transactions)

		summary := summaryData["summary"].(map[string]interface{})
		assert.Equal(t, float64(0), summary["total_transaction"])

		mockModel.AssertExpectations(t)
	})
}
