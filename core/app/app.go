package app

import (
	_ "fmt"
	"log"
	_ "os"
	_ "os/signal"
	_ "time"

	"online-store/core/app/handlers"
	"online-store/core/app/routes"
	"online-store/database" // Adjust import path as needed
	"online-store/package/config"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Server *gin.Engine
	DB     *sqlx.DB
}

func NewApp() *App {
	return &App{
		Server: gin.Default(),
	}
}

// Setup initializes the database and routes
func (a *App) Setup() {

	dbManager := database.InitConnection()
	db, err := dbManager.GetDB()
	if err != nil {
		log.Fatalf("Could not get database connection: %v", err)
	}
	a.DB = db

	// Initialize handler
	h := handlers.NewHandler(a.DB)

	// Setup routes
	routes.SetupRoutes(a.Server, h, a.DB)
}

// Run starts the HTTP server
func (a *App) Run() {
	cfg := config.GetConfig()
	go func() {
		if err := a.Server.Run(":" + cfg.AppPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Println("Server running on port:", cfg.AppPort)
}

// ShutDown gracefully shuts down the server
func (a *App) ShutDown() {
	if err := a.DB.Close(); err != nil {
		log.Println("Database connection could not be closed:", err)
	}
	log.Println("Server shutdown gracefully")
}
