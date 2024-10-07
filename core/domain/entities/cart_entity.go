package entities

import (
    "github.com/google/uuid"
    "time"
)

type Cart struct {
    ID         uuid.UUID  `db:"id" json:"id"`
    CustomerID uuid.UUID  `db:"customer_id" json:"customer_id"`
    CreatedAt  time.Time  `db:"created_at" json:"created_at"`
    UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
    Items   []ProductInCart  `json:"products,omitempty"` // Convenience field, not in DB
}


type CartProduct struct {
	CartID    uuid.UUID `db:"cart_id" json:"cart_id"`
	ProductID uuid.UUID `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
}

type ProductInCart struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
}