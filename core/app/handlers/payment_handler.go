package handlers

import (
	"context"
	"net/http"
	"online-store/core/domain/entities"
	"online-store/core/domain/repositories"
	"online-store/core/dto"
	_ "online-store/package/config"
	"online-store/package/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PaymentHandler struct {
	DB                *sqlx.DB
	InvoiceRepository repositories.InvoiceRepository
	PaymentRepository repositories.PaymentRepository
}

func NewPaymentHandler(db *sqlx.DB) *PaymentHandler {
	return &PaymentHandler{
		DB:                db,
		InvoiceRepository: repositories.NewInvoiceRepository(),
		PaymentRepository: repositories.NewPaymentRepository(),
	}
}

// CreatePayment creates a new payment
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var createPaymentDTO dto.CreatePaymentDTO
	if err := c.ShouldBindJSON(&createPaymentDTO); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	// Get customer ID from token
	userClaims, exists := c.Get("user")
	if !exists || userClaims == nil {
		response.RespondError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	customerID := userClaims.(entities.Customer).ID

	// Get invoice by ID
	invoice, err := h.InvoiceRepository.GetInvoiceByID(context.Background(), h.DB, createPaymentDTO.InvoiceID)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to get invoice", nil)
		return
	}

	// Check if invoice is already paid
	if invoice.Status == entities.InvoiceStatusPaid {
		response.RespondError(c, http.StatusBadRequest, "Invoice is already paid", nil)
		return
	}

	// Create payment

	invoiceID, err := uuid.Parse(createPaymentDTO.InvoiceID)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to parse invoice ID", nil)
		return
	}
	payment := entities.Payment{
		InvoiceID:  invoiceID,
		CustomerID: customerID,
		Amount:     invoice.Amount,
	}

	// Begin transaction
	tx, err := h.DB.BeginTxx(context.Background(), nil)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to start transaction", nil)
		return
	}

	// Create payment
	payment, err = h.PaymentRepository.SaveTx(context.Background(), tx, payment)
	if err != nil {
		tx.Rollback()
		response.RespondError(c, http.StatusInternalServerError, "Failed to create payment", nil)
		return
	}

	// Update invoice status
	invoice.Status = entities.InvoiceStatusPaid
	err = h.InvoiceRepository.FinishInvoice(context.Background(), tx, invoiceID.String(), payment.ID.String())
	if err != nil {
		tx.Rollback()
		response.RespondError(c, http.StatusInternalServerError, "Failed to update invoice", nil)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to commit transaction", nil)
		return
	}

	response.RespondSuccess(c, http.StatusOK, "Payment successful", payment)
}

// GetPaymentsByCustomerID gets all payments by customer ID
func (h *PaymentHandler) GetPaymentsByCustomerID(c *gin.Context) {
	// Get customer ID from token
	userClaims, exists := c.Get("user")
	if !exists || userClaims == nil {
		response.RespondError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	customerID := userClaims.(entities.Customer).ID

	// Get payments by customer ID
	payments, err := h.PaymentRepository.GetPaymentByCustomerID(context.Background(), h.DB, customerID.String())
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to get payments", nil)
		return
	}

	response.RespondSuccess(c, http.StatusOK, "Payments retrieved successfully", payments)
}
