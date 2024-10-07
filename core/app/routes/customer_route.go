package routes

import (
	"online-store/core/app/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// CustomerRoutes sets up the Customer-related routes.
func CustomerRoutes(router *gin.Engine, handler *handlers.CustomerHandler, db *sqlx.DB) {
	CustomerGroup := router.Group("/Customers")
	{   // Get Customers by category
		CustomerGroup.POST("SignUp", handler.CreateCustomer)          // Add a new Customer
		CustomerGroup.POST("SignIn", handler.SignInCustomer)          // Add a new Customer
	}
}
