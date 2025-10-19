package dto

import "payment-service/pkg/models"

type GetPaymentAttemptResponseDto struct {
	PaymentAttemptID string               `json:"payment_attempt_id"`
	OrderID          string               `json:"order_id"`
	PaymentInfoID    string               `json:"payment_info_id,omitempty"`
	Method           models.PaymentMethod `json:"method"`
	Status           models.PaymentStatus `json:"status"`
}
