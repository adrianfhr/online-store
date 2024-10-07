package entities

import (
	"github.com/google/uuid"
)

type Customer struct {
	ID       uuid.UUID `db:"id" json:"id"`         // Unique identifier
	Name     string    `db:"name" json:"name"`     // Customer's name
	Email    string    `db:"email" json:"email"`   // Customer's email (unique)
	Password string    `db:"password" json:"password"` // Customer's password
}
