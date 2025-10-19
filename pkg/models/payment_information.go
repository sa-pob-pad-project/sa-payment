package models

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
)

// PaymentMethod represents the payment method enum
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodPromptPay  PaymentMethod = "promptpay"
)

// Value implements the driver.Valuer interface
func (pm PaymentMethod) Value() (driver.Value, error) {
	return string(pm), nil
}

// Scan implements the sql.Scanner interface
func (pm *PaymentMethod) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	*pm = PaymentMethod(value.(string))
	return nil
}

// PaymentInformation represents the payment_informations table
type PaymentInformation struct {
	ID        uuid.UUID     `db:"id" json:"id"`
	UserID    uuid.UUID     `db:"user_id" json:"user_id"`
	Type      PaymentMethod `db:"type" json:"type"`
	Details   []byte        `db:"details" json:"details"`
	Version   int           `db:"version" json:"version"`
	CreatedAt time.Time     `db:"created_at" json:"created_at"`
}
