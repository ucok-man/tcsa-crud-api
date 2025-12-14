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

func TestCreateTransactionHandler(t *testing.T) {
	t.Run("successfully creates transaction", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		mockModel.On("Insert", mock.AnythingOfType("*data.Transaction")).
			Run(func(args mock.Arguments) {
				tx := args.Get(0).(*data.Transaction)
				tx.ID = 1
				tx.Version = 1
				tx.CreatedAt = time.Now()
				tx.UpdatedAt = time.Now()
			}).
			Return(nil)

		body := `{"user_id": 1, "amount": 10000}`
		ctx, rec := createTestContext(http.MethodPost, "/transactions", body)

		// Execute
		err := app.createTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		txData := response["data"].(map[string]interface{})
		assert.Equal(t, float64(1), txData["id"])
		assert.Equal(t, float64(1), txData["user_id"])
		assert.Equal(t, float64(10000), txData["amount"])
		assert.Equal(t, "pending", txData["status"])

		mockModel.AssertExpectations(t)
	})

	t.Run("returns validation error for missing user_id", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		body := `{"amount": 10000}`
		ctx, _ := createTestContext(http.MethodPost, "/transactions", body)

		// Execute
		err := app.createTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Insert")
	})

	t.Run("returns validation error for missing amount", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		body := `{"user_id": 1}`
		ctx, _ := createTestContext(http.MethodPost, "/transactions", body)

		// Execute
		err := app.createTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Insert")
	})

	t.Run("returns validation error for negative amount", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		body := `{"user_id": 1, "amount": -100}`
		ctx, _ := createTestContext(http.MethodPost, "/transactions", body)

		// Execute
		err := app.createTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Insert")
	})

	t.Run("returns validation error for negative user_id", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		body := `{"user_id": -1, "amount": 10000}`
		ctx, _ := createTestContext(http.MethodPost, "/transactions", body)

		// Execute
		err := app.createTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Insert")
	})

	t.Run("returns error for empty request body", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		body := `{}`
		ctx, _ := createTestContext(http.MethodPost, "/transactions", body)

		// Execute
		err := app.createTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Insert")
	})

	t.Run("returns error for malformed JSON", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		body := `{"user_id": 1, "amount": }`
		ctx, _ := createTestContext(http.MethodPost, "/transactions", body)

		// Execute
		err := app.createTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "Insert")
	})

	t.Run("returns error when database insert fails", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		mockModel.On("Insert", mock.AnythingOfType("*data.Transaction")).Return(assert.AnError)

		body := `{"user_id": 1, "amount": 10000}`
		ctx, _ := createTestContext(http.MethodPost, "/transactions", body)

		// Execute
		err := app.createTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertExpectations(t)
	})
}

func TestGetByIdTransactionHandler(t *testing.T) {
	t.Run("successfully gets transaction by id", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedTx := &data.Transaction{
			ID:        1,
			UserId:    1,
			Amount:    10000,
			Status:    data.TransactionStatusPending,
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockModel.On("GetById", 1).Return(expectedTx, nil)

		ctx, rec := createTestContext(http.MethodGet, "/transactions/1", "")
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Execute
		err := app.getByIdTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		txData := response["data"].(map[string]any)
		assert.Equal(t, float64(1), txData["id"])
		assert.Equal(t, float64(10000), txData["amount"])

		mockModel.AssertExpectations(t)
	})

	t.Run("returns 404 when transaction not found", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		mockModel.On("GetById", 999).Return(nil, data.ErrRecordNotFound)

		ctx, _ := createTestContext(http.MethodGet, "/transactions/999", "")
		ctx.SetParamNames("id")
		ctx.SetParamValues("999")

		// Execute
		err := app.getByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertExpectations(t)
	})

	t.Run("returns error for invalid id parameter", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/transactions/invalid", "")
		ctx.SetParamNames("id")
		ctx.SetParamValues("invalid")

		// Execute
		err := app.getByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetById")
	})

	t.Run("returns error when database getByid fails", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		mockModel.On("GetById", 1).Return(nil, assert.AnError)

		ctx, _ := createTestContext(http.MethodGet, "/transactions/1", "")
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Execute
		err := app.getByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertExpectations(t)
	})
}

func TestUpdateByIdTransactionHandler(t *testing.T) {
	t.Run("successfully updates transaction", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		existingTx := &data.Transaction{
			ID:        1,
			UserId:    1,
			Amount:    10000,
			Status:    data.TransactionStatusPending,
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockModel.On("GetById", 1).Return(existingTx, nil)
		mockModel.On("Update", mock.AnythingOfType("*data.Transaction")).Return(nil)

		body := `{"amount": 15000, "status": "success"}`
		ctx, rec := createTestContext(http.MethodPut, "/transactions/1", body)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Execute
		err := app.updateByIdTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		txData := response["data"].(map[string]any)
		assert.Equal(t, float64(15000), txData["amount"])
		assert.Equal(t, "success", txData["status"])

		mockModel.AssertExpectations(t)
	})

	t.Run("successfully updates only amount", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		existingTx := &data.Transaction{
			ID:        1,
			UserId:    1,
			Amount:    10000,
			Status:    data.TransactionStatusPending,
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockModel.On("GetById", 1).Return(existingTx, nil)
		mockModel.On("Update", mock.MatchedBy(func(tx *data.Transaction) bool {
			return tx.Amount == 20000 && tx.Status == data.TransactionStatusPending
		})).Return(nil)

		body := `{"amount": 20000}`
		ctx, rec := createTestContext(http.MethodPut, "/transactions/1", body)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Execute
		err := app.updateByIdTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		txData := response["data"].(map[string]any)
		assert.Equal(t, float64(20000), txData["amount"])
		assert.Equal(t, "pending", txData["status"]) // Status unchanged

		mockModel.AssertExpectations(t)
	})

	t.Run("successfully updates only status", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		existingTx := &data.Transaction{
			ID:        1,
			UserId:    1,
			Amount:    10000,
			Status:    data.TransactionStatusPending,
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockModel.On("GetById", 1).Return(existingTx, nil)
		mockModel.On("Update", mock.MatchedBy(func(tx *data.Transaction) bool {
			return tx.Amount == 10000 && tx.Status == data.TransactionStatusFailed
		})).Return(nil)

		body := `{"status": "failed"}`
		ctx, rec := createTestContext(http.MethodPut, "/transactions/1", body)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Execute
		err := app.updateByIdTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		txData := response["data"].(map[string]any)
		assert.Equal(t, float64(10000), txData["amount"]) // Amount unchanged
		assert.Equal(t, "failed", txData["status"])

		mockModel.AssertExpectations(t)
	})

	t.Run("returns 404 when transaction not found", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		mockModel.On("GetById", 999).Return(nil, data.ErrRecordNotFound)

		body := `{"amount": 15000}`
		ctx, _ := createTestContext(http.MethodPut, "/transactions/999", body)
		ctx.SetParamNames("id")
		ctx.SetParamValues("999")

		// Execute
		err := app.updateByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertExpectations(t)
		mockModel.AssertNotCalled(t, "Update")
	})

	t.Run("returns edit conflict error", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		existingTx := &data.Transaction{
			ID:      1,
			UserId:  1,
			Amount:  10000,
			Status:  data.TransactionStatusPending,
			Version: 1,
		}

		mockModel.On("GetById", 1).Return(existingTx, nil)
		mockModel.On("Update", mock.AnythingOfType("*data.Transaction")).Return(data.ErrEditConflict)

		body := `{"amount": 15000}`
		ctx, _ := createTestContext(http.MethodPut, "/transactions/1", body)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Execute
		err := app.updateByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertExpectations(t)
	})

	t.Run("returns validation error for invalid id parameter", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		body := `{"amount": 15000}`
		ctx, _ := createTestContext(http.MethodPut, "/transactions/invalid", body)
		ctx.SetParamNames("id")
		ctx.SetParamValues("invalid")

		// Execute
		err := app.updateByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetById")
		mockModel.AssertNotCalled(t, "Update")
	})

	t.Run("returns validation error for negative amount", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		body := `{"amount": -100}`
		ctx, _ := createTestContext(http.MethodPut, "/transactions/1", body)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Execute
		err := app.updateByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetById")
		mockModel.AssertNotCalled(t, "Update")
	})

	t.Run("returns validation error for invalid status", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		body := `{"status": "invalid_status"}`
		ctx, _ := createTestContext(http.MethodPut, "/transactions/1", body)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Execute
		err := app.updateByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetById")
		mockModel.AssertNotCalled(t, "Update")
	})

	t.Run("returns error for malformed JSON", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		body := `{"amount": "not_a_number"}`
		ctx, _ := createTestContext(http.MethodPut, "/transactions/1", body)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Execute
		err := app.updateByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetById")
		mockModel.AssertNotCalled(t, "Update")
	})

	t.Run("returns error for empty request body", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		body := `{}`
		ctx, rec := createTestContext(http.MethodPut, "/transactions/1", body)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		existingTx := &data.Transaction{
			ID:      1,
			UserId:  1,
			Amount:  10000,
			Status:  data.TransactionStatusPending,
			Version: 1,
		}

		mockModel.On("GetById", 1).Return(existingTx, nil)
		mockModel.On("Update", mock.AnythingOfType("*data.Transaction")).Return(nil)

		// Execute
		err := app.updateByIdTransactionHandler(ctx)

		// Assert - Should still succeed as both fields are optional
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockModel.AssertExpectations(t)
	})

	t.Run("returns error when database update fails", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		existingTx := &data.Transaction{
			ID:      1,
			UserId:  1,
			Amount:  10000,
			Status:  data.TransactionStatusPending,
			Version: 1,
		}

		mockModel.On("GetById", 1).Return(existingTx, nil)
		mockModel.On("Update", mock.AnythingOfType("*data.Transaction")).Return(assert.AnError)

		body := `{"amount": 15000}`
		ctx, _ := createTestContext(http.MethodPut, "/transactions/1", body)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Execute
		err := app.updateByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertExpectations(t)
	})
}

func TestDeleteTransactionHandler(t *testing.T) {
	t.Run("successfully deletes transaction", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		existingTx := &data.Transaction{
			ID:        1,
			UserId:    1,
			Amount:    10000,
			Status:    data.TransactionStatusPending,
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockModel.On("GetById", 1).Return(existingTx, nil)
		mockModel.On("DeleteOne", 1).Return(nil)

		ctx, rec := createTestContext(http.MethodDelete, "/transactions/1", "")
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Execute
		err := app.removeByIdTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockModel.AssertExpectations(t)
	})

	t.Run("returns validation error for invalid id parameter", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodDelete, "/transactions/invalid", "")
		ctx.SetParamNames("id")
		ctx.SetParamValues("invalid")

		// Execute
		err := app.removeByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetById")
		mockModel.AssertNotCalled(t, "DeleteOne")
	})

	t.Run("returns 404 when transaction not found", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		mockModel.On("GetById", 999).Return(nil, data.ErrRecordNotFound)

		ctx, _ := createTestContext(http.MethodDelete, "/transactions/999", "")
		ctx.SetParamNames("id")
		ctx.SetParamValues("999")

		// Execute
		err := app.removeByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertExpectations(t)
		mockModel.AssertNotCalled(t, "DeleteOne")
	})

	t.Run("returns error when database delete fails", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		existingTx := &data.Transaction{
			ID:      1,
			UserId:  1,
			Amount:  10000,
			Status:  data.TransactionStatusPending,
			Version: 1,
		}

		mockModel.On("GetById", 1).Return(existingTx, nil)
		mockModel.On("DeleteOne", 1).Return(assert.AnError)

		ctx, _ := createTestContext(http.MethodDelete, "/transactions/1", "")
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Execute
		err := app.removeByIdTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertExpectations(t)
	})
}

func TestGetAllTransactionHandler(t *testing.T) {
	t.Run("successfully gets all transactions with default pagination", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedTxs := []*data.Transaction{
			{ID: 1, UserId: 1, Amount: 10000, Status: data.TransactionStatusPending},
			{ID: 2, UserId: 2, Amount: 20000, Status: data.TransactionStatusSucces},
		}

		expectedMetadata := &data.Metadata{
			CurrentPage:  1,
			PageSize:     10,
			FirstPage:    1,
			LastPage:     1,
			TotalRecords: 2,
		}

		mockModel.On("GetAll", mock.AnythingOfType("data.TransactionGetAllParam")).
			Return(expectedTxs, expectedMetadata, nil)

		ctx, rec := createTestContext(http.MethodGet, "/transactions", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		txData := response["data"].([]any)
		assert.Len(t, txData, 2)

		metadata := response["metadata"].(map[string]any)
		assert.Equal(t, float64(1), metadata["current_page"])
		assert.Equal(t, float64(2), metadata["total_records"])

		mockModel.AssertExpectations(t)
	})

	t.Run("successfully filters by status", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedTxs := []*data.Transaction{
			{ID: 1, UserId: 1, Amount: 10000, Status: data.TransactionStatusSucces},
		}

		expectedMetadata := &data.Metadata{
			CurrentPage:  1,
			PageSize:     10,
			FirstPage:    1,
			LastPage:     1,
			TotalRecords: 1,
		}

		mockModel.On("GetAll", mock.MatchedBy(func(param data.TransactionGetAllParam) bool {
			return param.FilterStatus == "success"
		})).Return(expectedTxs, expectedMetadata, nil)

		ctx, rec := createTestContext(http.MethodGet, "/transactions?status=success", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockModel.AssertExpectations(t)
	})

	t.Run("successfully filters by user_id", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedTxs := []*data.Transaction{
			{ID: 1, UserId: 5, Amount: 10000, Status: data.TransactionStatusPending},
			{ID: 2, UserId: 5, Amount: 15000, Status: data.TransactionStatusSucces},
		}

		expectedMetadata := &data.Metadata{
			CurrentPage:  1,
			PageSize:     10,
			FirstPage:    1,
			LastPage:     1,
			TotalRecords: 2,
		}

		mockModel.On("GetAll", mock.MatchedBy(func(param data.TransactionGetAllParam) bool {
			return param.FilterUserId == 5
		})).Return(expectedTxs, expectedMetadata, nil)

		ctx, rec := createTestContext(http.MethodGet, "/transactions?user_id=5", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		txData := response["data"].([]any)
		assert.Len(t, txData, 2)

		mockModel.AssertExpectations(t)
	})

	t.Run("successfully applies custom pagination", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedTxs := []*data.Transaction{
			{ID: 21, UserId: 1, Amount: 10000, Status: data.TransactionStatusPending},
		}

		expectedMetadata := &data.Metadata{
			CurrentPage:  3,
			PageSize:     20,
			FirstPage:    1,
			LastPage:     5,
			TotalRecords: 100,
		}

		mockModel.On("GetAll", mock.MatchedBy(func(param data.TransactionGetAllParam) bool {
			return param.Page == 3 && param.PageSize == 20
		})).Return(expectedTxs, expectedMetadata, nil)

		ctx, rec := createTestContext(http.MethodGet, "/transactions?page=3&page_size=20", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response envelope
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse JSON response")

		metadata := response["metadata"].(map[string]any)
		assert.Equal(t, float64(3), metadata["current_page"])
		assert.Equal(t, float64(20), metadata["page_size"])

		mockModel.AssertExpectations(t)
	})

	t.Run("successfully applies custom sorting", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedTxs := []*data.Transaction{
			{ID: 2, UserId: 1, Amount: 20000, Status: data.TransactionStatusSucces},
			{ID: 1, UserId: 1, Amount: 10000, Status: data.TransactionStatusPending},
		}

		expectedMetadata := &data.Metadata{
			CurrentPage:  1,
			PageSize:     10,
			FirstPage:    1,
			LastPage:     1,
			TotalRecords: 2,
		}

		mockModel.On("GetAll", mock.MatchedBy(func(param data.TransactionGetAllParam) bool {
			return param.SortColumn == "amount" && param.SortDirection == "DESC"
		})).Return(expectedTxs, expectedMetadata, nil)

		ctx, rec := createTestContext(http.MethodGet, "/transactions?sort_by=-amount", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockModel.AssertExpectations(t)
	})

	t.Run("successfully combines multiple filters", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		expectedTxs := []*data.Transaction{
			{ID: 1, UserId: 3, Amount: 10000, Status: data.TransactionStatusSucces},
		}

		expectedMetadata := &data.Metadata{
			CurrentPage:  1,
			PageSize:     10,
			FirstPage:    1,
			LastPage:     1,
			TotalRecords: 1,
		}

		mockModel.On("GetAll", mock.MatchedBy(func(param data.TransactionGetAllParam) bool {
			return param.FilterStatus == "success" &&
				param.FilterUserId == 3 &&
				param.SortColumn == "created_at" &&
				param.SortDirection == "DESC"
		})).Return(expectedTxs, expectedMetadata, nil)

		ctx, rec := createTestContext(http.MethodGet, "/transactions?status=success&user_id=3&sort_by=-created_at", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockModel.AssertExpectations(t)
	})

	t.Run("returns validation error for invalid page number", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/transactions?page=0", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetAll")
	})

	t.Run("returns validation error for page exceeding maximum", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/transactions?page=1001", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetAll")
	})

	t.Run("returns validation error for invalid page size", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/transactions?page_size=0", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetAll")
	})

	t.Run("returns validation error for page size exceeding maximum", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/transactions?page_size=101", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetAll")
	})

	t.Run("returns validation error for invalid sort_by value", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/transactions?sort_by=invalid_field", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetAll")
	})

	t.Run("returns validation error for invalid status filter", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/transactions?status=invalid_status", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetAll")
	})

	t.Run("returns validation error for invalid user_id", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		ctx, _ := createTestContext(http.MethodGet, "/transactions?user_id=0", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertNotCalled(t, "GetAll")
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		// Setup
		mockModel := new(data.MockTransactionModel)
		app := createTestApp(t, data.Models{Transactions: mockModel})

		mockModel.On("GetAll", mock.AnythingOfType("data.TransactionGetAllParam")).
			Return(nil, nil, assert.AnError)

		ctx, _ := createTestContext(http.MethodGet, "/transactions", "")

		// Execute
		err := app.getAllTransactionHandler(ctx)

		// Assert
		assert.Error(t, err)
		mockModel.AssertExpectations(t)
	})
}
