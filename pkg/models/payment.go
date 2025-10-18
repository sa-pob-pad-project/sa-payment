package models

import (
	"time"

	"github.com/google/uuid"
)

// Payment represents the payments table
type Payment struct {
	ID        uuid.UUID `db:"id" json:"id"`
	AttemptID uuid.UUID `db:"attempt_id" json:"attempt_id"`
	Amount    float64   `db:"amount" json:"amount"`
	OrderID   uuid.UUID `db:"order_id" json:"order_id"`
	PaidAt    time.Time `db:"paid_at" json:"paid_at"`
}
