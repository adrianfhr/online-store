package entities

import (
	"github.com/google/uuid"
)

type CartItem struct {
	ID         uuid.UUID `db:"id" json:"id"` // Unique identifier
	CustomerID uuid.UUID `db:"customer_id" json:"customer_id"` // Customer's ID
	ProductID  uuid.UUID `db:"product_id" json:"product_id"`   // Product's ID
	Quantity   int       `db:"quantity" json:"quantity"`       // Quantity of the product
}
