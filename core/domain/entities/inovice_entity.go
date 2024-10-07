package entities

import (
	"github.com/google/uuid"
	"time"
)

type Invoice struct {
	ID         uuid.UUID     `db:"id" json:"id"`
	CustomerID uuid.UUID     `db:"customer_id" json:"customer_id"`
	Amount     float64       `db:"amount" json:"amount"`
	Status     string        `db:"status" json:"status"` // Pending, Paid, Expired
	CreatedAt  time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at" json:"updated_at"`
	Expiration time.Time     `db:"expiration" json:"expiration"`
	Items      []InvoiceItem `json:"items,omitempty"` // Convenience field, not in DB
}

type InvoiceItem struct {
	ID          uuid.UUID `db:"id" json:"id"`
	InvoiceID   uuid.UUID `db:"invoice_id" json:"invoice_id"`
	ProductID   uuid.UUID `db:"product_id" json:"product_id"`
	ProductName string    `db:"product_name" json:"product_name"`
	Quantity    int       `db:"quantity" json:"quantity"`
	Price       float64   `db:"price" json:"price"`
}

// invoice status pending, paid, expired
const (
	InvoiceStatusPending = "pending"
	InvoiceStatusPaid    = "paid"
	InvoiceStatusExpired = "expired"
)