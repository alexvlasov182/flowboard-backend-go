package main

import (
	"flowboard-backend-go/internal/database"
	"flowboard-backend-go/internal/pages"
	"flowboard-backend-go/internal/users"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load ENV
	err := godotenv.Load()
	if err != nil {

		log.Printf("No .env file found or error loading it: %v", err)
	}

	// Connect to database
	db := database.Connect()

	// Auto-migrate database schemas
	err = db.AutoMigrate(&users.User{}, &pages.Page{})
	if err != nil {
		log.Fatal("Failed to auto-migrate database schemas:", err)
	}

	log.Println("Database migrated successfully")
	log.Println("FlowBoard Go API ready!")
}
