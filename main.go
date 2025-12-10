package main

import (
	"log"
	"os"
	"warehouse/config"
	"warehouse/routers"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.DBInit()
	config.RunMigration()

	r := routers.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server running on port %s", port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
