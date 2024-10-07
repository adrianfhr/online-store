package entities

import (
    "time"
    "github.com/google/uuid"
)

type Payment struct {
    ID           uuid.UUID  `db:"id" json:"id"`
    InvoiceID    uuid.UUID  `db:"invoice_id" json:"invoice_id"`
	CustomerID   uuid.UUID  `db:"customer_id" json:"customer_id"`
    Amount       float64    `db:"amount" json:"amount"`
    Status       string     `db:"status" json:"status"` // pending, completed, failed
    PaymentDate  time.Time  `db:"payment_date" json:"payment_date"`
    CreatedAt    time.Time  `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

const (
	PaymentStatusPending   = "pending"
	PaymentStatusCompleted = "completed"
	PaymentStatusFailed    = "failed"
)