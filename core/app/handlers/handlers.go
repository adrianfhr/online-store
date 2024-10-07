package handlers

import (
	"github.com/jmoiron/sqlx"
)

// Handler represents a struct that holds references to various handlers.
type Handler struct {
	ProductHandler *ProductHandler
	CustomerHandler *CustomerHandler
}

// NewHandler initializes the main handler and injects dependencies.
func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		ProductHandler: NewProductHandler(db),
		CustomerHandler: NewCustomerHandler(db),
	}
}
