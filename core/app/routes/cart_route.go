package routes

import (
	"online-store/core/app/handlers"
	"online-store/core/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// CartRoutes sets up the Cart-related routes.
func CartRoutes(router *gin.Engine, handler *handlers.CartHandler,   db *sqlx.DB) {
	CartGroup := router.Group("/Cart")
	{ // Get Carts by category
		CartGroup.POST("Items", middleware.RequireAuthMiddleware(db), handler.AddToCart)            // Add a new Cart
		CartGroup.GET("Items", middleware.RequireAuthMiddleware(db), handler.GetCartWithProducts)   // Add a new Cart
		CartGroup.DELETE("Items", middleware.RequireAuthMiddleware(db), handler.RemoveItemFromCart) // Add a new Cart
		CartGroup.POST("Checkout", middleware.RequireAuthMiddleware(db), handler.CreateInvoice)           // Add a new Cart
	}
}
