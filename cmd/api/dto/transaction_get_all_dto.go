package dto

type TransactionGetAllDTO struct {
	Pagination struct {
		Page     int `validate:"omitempty,min=1,max=1000"`
		PageSize int `validate:"omitempty,min=1,max=100"`
		Limit    *int
		Offset   *int
	}
	Sort struct {
		RawValue    string `validate:"omitempty,oneof=amount status created_at -amount -status -created_at ''"`
		ColumnValue *string
		Direction   *string
	}
	Filter struct {
		Status string `validate:"omitempty,oneof=pending failed success ''"`
	}
}
