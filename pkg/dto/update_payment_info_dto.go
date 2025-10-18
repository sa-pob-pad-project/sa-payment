package dto

import (
	"encoding/json"
	"time"

	"order-service/pkg/models"
)

type UpdatePaymentInfoRequestDto struct {
	ID            string               `json:"id" validate:"required"`
	PaymentMethod models.PaymentMethod `json:"payment_method" validate:"required,oneof=credit_card promptpay"`
	Details       json.RawMessage      `json:"details" validate:"required"`
}

type UpdatePaymentInfoResponseDto struct {
	ID            string               `json:"id"`
	UserID        string               `json:"user_id"`
	PaymentMethod models.PaymentMethod `json:"payment_method"`
	Details       json.RawMessage      `json:"details"`
	Version       int                  `json:"version"`
	CreatedAt     time.Time            `json:"created_at"`
	// UpdateAt      time.Time            `json:"update_at"`
}
