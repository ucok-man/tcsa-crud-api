package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ucok-man/tcsa/internal/utility"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusFailed  TransactionStatus = "failed"
	TransactionStatusSucces  TransactionStatus = "success"
)

type Transaction struct {
	ID        int               `json:"id"`
	UserId    int               `json:"user_id"`
	Amount    int               `json:"amount"`
	Status    TransactionStatus `json:"status"`
	Version   int               `json:"-"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Summary struct {
	CountTotal int `json:"total_transaction"`
	Pending    struct {
		Count          int     `json:"count"`
		RatePercentage float64 `json:"rate_percentage"`
	} `json:"pending"`
	Success struct {
		Count          int     `json:"count"`
		RatePercentage float64 `json:"rate_percentage"`
	} `json:"success"`
	Failed struct {
		Count          int     `json:"count"`
		RatePercentage float64 `json:"rate_percentage"`
	} `json:"failed"`
}

type TransactionSummary struct {
	Transactions []*Transaction `json:"transactions"`
	Summary      Summary        `json:"summary"`
}

type TransactionModel struct {
	db *sql.DB
}

func (m TransactionModel) Insert(transaction *Transaction) error {
	query := `
        INSERT INTO transactions (user_id, amount, status)
        VALUES ($1, $2, $3)
        RETURNING id, version, created_at, updated_at`
	args := []any{transaction.UserId, transaction.Amount, transaction.Status}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.db.QueryRowContext(ctx, query, args...).Scan(
		&transaction.ID,
		&transaction.Version,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
}

type TransactionGetAllParam struct {
	Page          int
	PageSize      int
	PageOffset    int
	SortColumn    string
	SortDirection string
	FilterStatus  string
	FilterUserId  int
}

func (m TransactionModel) GetAll(param TransactionGetAllParam) ([]*Transaction, *Metadata, error) {
	query := fmt.Sprintf(`
	    SELECT 
			count(*) OVER() as total_count, 
			id, user_id, amount, status, version, created_at, updated_at
	    FROM transactions
	    WHERE 
			(CASE 
				WHEN $1 = '' THEN TRUE
				ELSE status = $1
			END)
			AND
			(CASE 
				WHEN $2 = 0 THEN TRUE
				ELSE user_id = $2
			END)		
	    ORDER BY %s %s, id ASC
	    LIMIT $3 OFFSET $4`, param.SortColumn, param.SortDirection,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{param.FilterStatus, param.FilterUserId, param.PageSize, param.PageOffset}

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var totalRecords int
	var transactions []*Transaction

	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&totalRecords, // count from window function
			&transaction.ID,
			&transaction.UserId,
			&transaction.Amount,
			&transaction.Status,
			&transaction.Version,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}

		transactions = append(transactions, &transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	metadata := calculateMetadata(totalRecords, param.Page, param.PageSize)
	return transactions, &metadata, nil
}

func (m TransactionModel) GetById(id int) (*Transaction, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, user_id, amount, status, version, created_at, updated_at
		FROM transactions
		WHERE id = $1`

	var transaction Transaction

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.UserId,
		&transaction.Amount,
		&transaction.Status,
		&transaction.Version,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &transaction, nil
}

func (m TransactionModel) Update(transaction *Transaction) error {
	query := `
        UPDATE transactions
        SET amount = $1, status = $2, updated_at=$3, version = version + 1
        WHERE id = $4 AND version = $5
        RETURNING version`

	args := []any{
		&transaction.Amount,
		&transaction.Status,
		&transaction.UpdatedAt,
		&transaction.ID,
		&transaction.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, args...).Scan(&transaction.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m TransactionModel) DeleteOne(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM transactions WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

type TransactionSummaryParam struct {
	Page            int
	PageSize        int
	PageOffset      int
	SortColumn      string
	SortDirection   string
	FilterDateRange int
	FilterUserId    int
}

func (m TransactionModel) Summary(param TransactionSummaryParam) (*TransactionSummary, *Metadata, error) {
	query := fmt.Sprintf(`
    SELECT 
        count(*) OVER() as total_count,
        count(*) FILTER (WHERE status = 'success') OVER() as success_count,
        count(*) FILTER (WHERE status = 'pending') OVER() as pending_count,
        count(*) FILTER (WHERE status = 'failed') OVER() as failed_count,
        id, user_id, amount, status, version, created_at, updated_at
    FROM transactions
    WHERE
        (CASE 
            WHEN $1 = 0 THEN TRUE
            ELSE created_at >= CURRENT_DATE - INTERVAL '1 day' * ($1 - 1)
        END) 
        AND
        (CASE 
            WHEN $2 = 0 THEN TRUE
            ELSE user_id = $2
        END)
    ORDER BY %s %s, id ASC
    LIMIT $3 OFFSET $4`, param.SortColumn, param.SortDirection,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{param.FilterDateRange, param.FilterUserId, param.PageSize, param.PageOffset}

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var summary Summary
	var transactions []*Transaction

	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&summary.CountTotal,
			&summary.Success.Count,
			&summary.Pending.Count,
			&summary.Failed.Count,
			&transaction.ID,
			&transaction.UserId,
			&transaction.Amount,
			&transaction.Status,
			&transaction.Version,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}

		transactions = append(transactions, &transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	if summary.CountTotal > 0 {
		summary.Success.RatePercentage = utility.Round2(
			float64(summary.Success.Count) / float64(summary.CountTotal) * 100,
		)

		summary.Pending.RatePercentage = utility.Round2(
			float64(summary.Pending.Count) / float64(summary.CountTotal) * 100,
		)

		summary.Failed.RatePercentage = utility.Round2(
			float64(summary.Failed.Count) / float64(summary.CountTotal) * 100,
		)
	}

	metadata := calculateMetadata(summary.CountTotal, param.Page, param.PageSize)

	return &TransactionSummary{
		Transactions: transactions,
		Summary:      summary,
	}, &metadata, err
}
