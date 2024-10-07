package entities

import (
	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID `db:"id" json:"id"` // Unique identifier
	CustomerID uuid.UUID `db:"customer_id" json:"customer_id"` // Customer's ID
	Total      float64   `db:"total" json:"total"` // Total amount of the order
	Status     string    `db:"status" json:"status"` // Status of the order
	CreatedAt  uint      `db:"created_at" json:"created_at"` // Time the order was created
}