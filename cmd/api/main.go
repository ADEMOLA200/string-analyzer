package main

import (
	"log"

	"github.com/ADEMOLA200/string-analyzer/cmd/internal/routes"
	"github.com/ADEMOLA200/string-analyzer/config"
)

func main() {
	cfg := config.Load()

	router := routes.SetupRouter()

	log.Printf("Server starting on port %s in %s environment", cfg.Port, cfg.Environment)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
