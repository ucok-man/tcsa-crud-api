package dto

type TransactionCreateDTO struct {
	UserId int `json:"user_id" validate:"required,min=1"`
	Amount int `json:"amount" validate:"required,min=1"`
}
