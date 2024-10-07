package entities

import (
	"github.com/google/uuid"
)

type Product struct {
	ID       uuid.UUID `db:"id" json:"id"`         // Unique identifier
	Name     string    `db:"name" json:"name"`     // Product's name
	Category string    `db:"category" json:"category"` // Product's category
	Price    float64   `db:"price" json:"price"`   // Product's price
}