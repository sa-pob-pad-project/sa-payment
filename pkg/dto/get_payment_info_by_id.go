package dto

import (
	"encoding/json"
	"order-service/pkg/models"
	"time"
)

type PaymentInfoDto struct {
	ID            string               `json:"id"`
	UserID        string               `json:"user_id"`
	PaymentMethod models.PaymentMethod `json:"payment_method"`
	Details       json.RawMessage      `json:"details"`
	Version       int                  `json:"version"`
	CreatedAt     time.Time            `json:"created_at"`
}

type GetPaymentInfoByIDResponseDto struct {
	PaymentInfo PaymentInfoDto `json:"payment_info"`
}

// Conversion functions
func ToPaymentInfoDto(info *models.PaymentInformation) PaymentInfoDto {
	return PaymentInfoDto{
		ID:            info.ID.String(),
		UserID:        info.UserID.String(),
		PaymentMethod: info.Type,
		Details:       info.Details,
		Version:       info.Version,
		CreatedAt:     info.CreatedAt,
	}
}
