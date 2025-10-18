package dto

type CreatePaymentRequestDto struct {
	PaymentAttemptID string  `json:"payment_attempt_id" validate:"required"`
	Amount           float64 `json:"amount" validate:"required,gte=0"`
}

type CreatePaymentResponseDto struct {
	PaymentID string  `json:"payment_id"`
	AttemptID string  `json:"attempt_id"`
	OrderID   string  `json:"order_id"`
	Amount    float64 `json:"amount"`
	PaidAt    string  `json:"paid_at"`
}
