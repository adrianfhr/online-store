package routes

import (
	"online-store/core/app/handlers"
	"online-store/core/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// ProductRoutes sets up the product-related routes.
func ProductRoutes(router *gin.Engine, handler *handlers.ProductHandler, db *sqlx.DB) {
	productGroup := router.Group("/Products")
	{
		productGroup.GET("", middleware.RequireAuthMiddleware(db), handler.GetProducts)          // Get products by category
		productGroup.POST("", middleware.RequireAuthMiddleware(db),handler.AddProduct)          // Add a new product
		productGroup.GET("/Categories", middleware.RequireAuthMiddleware(db), handler.GetCategories) // Get all categories
	}
}
