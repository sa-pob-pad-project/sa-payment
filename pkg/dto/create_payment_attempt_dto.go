package dto

type CreatePaymentAttemptRequestDto struct {
	OrderID       string `json:"order_id" validate:"required"`
	PaymentInfoID string `json:"payment_info_id" validate:"required"`
}

type CreatePaymentAttemptResponseDto struct {
	PaymentAttemptID string `json:"payment_attempt_id"`
}
