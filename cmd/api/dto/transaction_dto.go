package dto

type TransactionCreateDTO struct {
	UserId int `json:"user_id" validate:"required,min=1"`
	Amount int `json:"amount" validate:"required,min=1"`
}

type TransactionUpdateDTO struct {
	TransactionId int     `param:"id" validate:"required,min=1"`
	Amount        *int    `json:"amount" validate:"omitempty,min=1"`
	Status        *string `json:"status" validate:"omitempty,oneof=pending failed success"`
}

type TransactionGetAllDTO struct {
	Pagination struct {
		Page     *int `query:"page" validate:"omitempty,min=1,max=1000"`
		PageSize *int `query:"page_size" validate:"omitempty,min=1,max=100"`
	}
	Sort struct {
		Value *string `query:"sort_by" validate:"omitempty,oneof=id user_id amount status created_at -id -user_id -amount -status -created_at"`
	}
	Filter struct {
		Status *string `query:"status" validate:"omitempty,oneof=pending failed success"`
		UserId *int    `query:"user_id" validate:"omitempty,min=1"`
	}
}

type TransactionParamIdDTO struct {
	TransactionId int `param:"id" validate:"required,min=1"`
}

type TransactionSummaryDTO struct {
	Pagination struct {
		Page     *int `query:"page" validate:"omitempty,min=1,max=1000"`
		PageSize *int `query:"page_size" validate:"omitempty,min=1,max=100"`
	}
	Sort struct {
		Value *string `query:"sort_by" validate:"omitempty,oneof=id user_id amount status created_at -id -user_id -amount -status -created_at"`
	}
	Filter struct {
		DateRange *int `query:"date_range" validate:"omitempty,min=1,max=366"`
		UserId    *int `query:"user_id" validate:"omitempty,min=1"`
	}
}
