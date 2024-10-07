package routes

import (
	"online-store/core/app/handlers"
	"online-store/core/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// InvoiceRoutes sets up the Invoice-related routes.
func InvoiceRoutes(router *gin.Engine, handler *handlers.InvoiceHandler, db *sqlx.DB) {
	InvoiceGroup := router.Group("/Invoices")
	{ // Get Invoices by category
		InvoiceGroup.GET("", middleware.RequireAuthMiddleware(db), handler.GetInvoicesByCustomerID) // Get Invoices by customer ID
	}
}
