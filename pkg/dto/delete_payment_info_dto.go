package dto

type DeletePaymentInfoRequestDto struct {
	ID string `json:"id" validate:"required"`
}

type DeletePaymentInfoResponseDto struct {
	ID        string `json:"id"`
	DeletedAt string `json:"deleted_at"`
}
