package routes

import (
	"net/http"
	"online-store/core/app/handlers"
	"online-store/package/response"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// SetupRoutes initializes all application routes.
func SetupRoutes(router *gin.Engine, handler *handlers.Handler, db *sqlx.DB) {
	// Register product routes
	// Use custom 404 handler
	router.NoRoute(func(c *gin.Context) {
		response.RespondError(c, http.StatusNotFound, "Resource not found", nil)
	})
	
	ProductRoutes(router, handler.ProductHandler, db)
	CustomerRoutes(router, handler.CustomerHandler, db)

	// You can add more routes here for other entities, e.g., auth routes, order routes, etc.
}
