package dto

import (
	"payment-service/pkg/models"
)

type CreditCardDetails struct {
	CardNumber     string `json:"card_number" validate:"required"`
	Cvv            string `json:"cvv" validate:"required"`
	ExpiryMonth    int    `json:"expiry_month" validate:"required,min=1,max=12"`
	ExpiryYear     int    `json:"expiry_year" validate:"required"`
	CardHolderName string `json:"card_holder_name" validate:"required"`
}

type PromptPayDetails struct {
	PromptPayID   string `json:"promptpay_id" validate:"required"`
	PromptPayType string `json:"promptpay_type" validate:"required"`
}

// Request
type CreatePaymentInfoRequestDto struct {
	PaymentMethod models.PaymentMethod `json:"payment_method" validate:"required,oneof=credit_card promptpay"`
	Details       []byte               `json:"details" validate:"required"`
}

// Response
//
//	type CreatePaymentInfoResponseDto struct {
//		ID            string               `json:"id"`
//		UserID        string               `json:"user_id"`
//		PaymentMethod models.PaymentMethod `json:"payment_method"`
//		Details       []byte      `json:"details"`
//		Version       int                  `json:"version"`
//		CreatedAt     time.Time            `json:"created_at"`
//	}
type CreatePaymentInfoResponseDto struct {
	PaymentInfo PaymentInfoDto `json:"payment_info"`
}
