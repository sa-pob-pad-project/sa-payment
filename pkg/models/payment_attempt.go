package models

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the payment_status enum
type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
)

// Value implements the driver.Valuer interface
func (ps PaymentStatus) Value() (driver.Value, error) {
	return string(ps), nil
}

// Scan implements the sql.Scanner interface
func (ps *PaymentStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	*ps = PaymentStatus(value.(string))
	return nil
}

// PaymentAttempt represents the payment_attempts table
type PaymentAttempt struct {
	ID                   uuid.UUID     `db:"id" json:"id"`
	OrderID              uuid.UUID     `db:"order_id" json:"order_id"`
	PaymentInformationID *uuid.UUID    `db:"payment_information_id" json:"payment_information_id"`
	Method               PaymentMethod `db:"method" json:"method"`
	Status               PaymentStatus `db:"status" json:"status"`
	CreatedAt            time.Time     `db:"created_at" json:"created_at"`
}
