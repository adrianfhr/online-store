package main

import (
	"fmt"
	_ "online-store/package/config"
	"os"
	"os/signal"

	"online-store/core/app"

	_ "github.com/golang-migrate/migrate"
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

	// log.GetLogger().Info("main", fmt.Sprintf("Server %s stopped", config.GetConfig().AppName), "gracefull", "")
	fmt.Println("Server stopped")
}
