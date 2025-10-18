package dto

import "order-service/pkg/models"

type UpdatePaymentAttemptRequestDto struct {
	PaymentAttemptID string               `json:"payment_attempt_id" validate:"required"`
	Status           models.PaymentStatus `json:"status" validate:"required"`
}

type UpdatePaymentAttemptResponseDto struct {
	PaymentAttemptID string               `json:"payment_attempt_id"`
	OrderID          string               `json:"order_id"`
	PaymentInfoID    string               `json:"payment_info_id,omitempty"`
	Method           models.PaymentMethod `json:"method"`
	Status           models.PaymentStatus `json:"status"`
}
