package main

import (
	"fmt"
	_ "online-store/package/config"
	"os"
	"os/signal"

	"online-store/core/app"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

func main() {
	apps := app.NewApp()

	apps.Setup()

	// start server
	apps.Run()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	apps.ShutDown()

	// Menjalankan migrasi
    m, err := migrate.New("file://./database/migrations", "postgres://postgres:ad681789@localhost:5432/onlinestore?sslmode=disable")
    if err != nil {
        fmt.Println("Error creating migration instance: ", err)
    }

    // Menjalankan migrasi
    if err := m.Up(); err != nil {
        if err != migrate.ErrNoChange {
            fmt.Println("Error running migration: ", err)
        } else {
            fmt.Println("No migration was applied")
        }
    }

	// log.GetLogger().Info("main", fmt.Sprintf("Server %s stopped", config.GetConfig().AppName), "gracefull", "")
	fmt.Println("Server stopped")
}
