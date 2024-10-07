package handlers

import (
	"context"
	"net/http"
	"online-store/core/domain/entities"
	"online-store/core/domain/repositories"
	_ "online-store/core/dto"
	_ "online-store/package/config"
	"online-store/package/response"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type InvoiceHandler struct {
	DB                *sqlx.DB
	CartRepository    repositories.CartRepository
	ProductRepository repositories.ProductRepository
	InvoiceRepository repositories.InvoiceRepository
}

// NewCartHandler creates a new CartHandler
func NewInvoiceHandler(db *sqlx.DB) *InvoiceHandler {
	return &InvoiceHandler{
		DB:                db,
		CartRepository:    repositories.NewCartRepository(),
		ProductRepository: repositories.NewProductRepository(),
		InvoiceRepository: repositories.NewInvoiceRepository(),
	}
}

// CreateInvoice creates a new invoice
// GetInvoicesByCustomerID retrieves all invoices for a specific customer
func (h *InvoiceHandler) GetInvoicesByCustomerID(c *gin.Context) {
	// Get customer ID from token
	userClaims, exists := c.Get("user")
	if !exists || userClaims == nil {
		response.RespondError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	customerID := userClaims.(entities.Customer).ID

	invoices, err := h.InvoiceRepository.GetInvoicesByCustomerID(context.Background(), h.DB, customerID.String())
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to get invoices", nil)
		return
	}

	response.RespondSuccess(c, http.StatusOK, response.GetSuccess(), invoices)
}