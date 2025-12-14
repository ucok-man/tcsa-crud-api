package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type TransactionModeler interface {
	Insert(transaction *Transaction) error
	GetAll(param TransactionGetAllParam) ([]*Transaction, *Metadata, error)
	GetById(id int) (*Transaction, error)
	Update(transaction *Transaction) error
	DeleteOne(id int) error
	Summary(param TransactionSummaryParam) (*TransactionSummary, *Metadata, error)
}

type Models struct {
	Transactions TransactionModeler
}

func NewModels(db *sql.DB) Models {
	return Models{
		Transactions: TransactionModel{db: db},
	}
}
