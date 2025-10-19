package dto

import (
	"time"

	"payment-service/pkg/models"
)

type PaymentDto struct {
	PaymentID string  `json:"payment_id"`
	AttemptID string  `json:"attempt_id"`
	OrderID   string  `json:"order_id"`
	Amount    float64 `json:"amount"`
	PaidAt    string  `json:"paid_at"`
}

type GetAllPaymentsResponseDto struct {
	Payments []PaymentDto `json:"payments"`
}

type GetPaymentByIDResponseDto struct {
	Payment PaymentDto `json:"payment"`
}

func ToPaymentDto(payment *models.Payment) PaymentDto {
	return PaymentDto{
		PaymentID: payment.ID.String(),
		AttemptID: payment.AttemptID.String(),
		OrderID:   payment.OrderID.String(),
		Amount:    payment.Amount,
		PaidAt:    payment.PaidAt.Format(time.RFC3339),
	}
}

func ToPaymentDtoList(payments []models.Payment) []PaymentDto {
	result := make([]PaymentDto, len(payments))
	for i := range payments {
		result[i] = ToPaymentDto(&payments[i])
	}
	return result
}
