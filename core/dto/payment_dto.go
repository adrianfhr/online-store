package dto

type CreatePaymentDTO struct {
	InvoiceID string  `json:"invoice_id" binding:"required"`
}