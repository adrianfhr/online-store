package routes

import (
	"online-store/core/app/handlers"
	"online-store/core/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PaymentRoutes sets up the Payment-related routes.
func PaymentRoutes(router *gin.Engine, handler *handlers.PaymentHandler, db *sqlx.DB) {
	PaymentGroup := router.Group("/Payments")
	{ // Get Payments by category
		PaymentGroup.GET("", middleware.RequireAuthMiddleware(db), handler.GetPaymentsByCustomerID) // Get Payments by customer ID
		PaymentGroup.POST("", middleware.RequireAuthMiddleware(db), handler.CreatePayment)          // Get all Payments
	}
}
