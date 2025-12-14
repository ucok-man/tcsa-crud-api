package data

import (
	"github.com/stretchr/testify/mock"
)

type MockTransactionModel struct {
	mock.Mock
}

func (m *MockTransactionModel) Insert(transaction *Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionModel) GetAll(param TransactionGetAllParam) ([]*Transaction, *Metadata, error) {
	args := m.Called(param)

	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil, args.Error(2)
		}
		return nil, args.Get(1).(*Metadata), args.Error(2)
	}

	return args.Get(0).([]*Transaction), args.Get(1).(*Metadata), args.Error(2)
}

func (m *MockTransactionModel) GetById(id int) (*Transaction, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Transaction), args.Error(1)
}

func (m *MockTransactionModel) Update(transaction *Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionModel) DeleteOne(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTransactionModel) Summary(param TransactionSummaryParam) (*TransactionSummary, *Metadata, error) {
	args := m.Called(param)

	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil, args.Error(2)
		}
		return nil, args.Get(1).(*Metadata), args.Error(2)
	}

	return args.Get(0).(*TransactionSummary), args.Get(1).(*Metadata), args.Error(2)
}
