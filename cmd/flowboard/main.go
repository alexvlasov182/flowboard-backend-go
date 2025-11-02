package main

import (
	"flowboard-backend-go/internal/database"
	"flowboard-backend-go/internal/middleware"
	"flowboard-backend-go/internal/pages"
	_users "flowboard-backend-go/internal/users"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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
	err = db.AutoMigrate(&_users.User{}, &pages.Page{})
	if err != nil {
		log.Fatal("Failed to auto-migrate database schemas:", err)
	}

	log.Println("Database migrated successfully")
	log.Println("FlowBoard Go API ready!")

	// create layers
	repo := _users.NewRepository(db)
	service := _users.NewService(repo)
	jwtMgr := middleware.NewJWTManager()
	handler := _users.NewHandler(service, jwtMgr)

	r := gin.Default()
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", handler.Register)
			auth.POST("/login", handler.Login)
		}

		// protected user routes
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtMgr))
		{
			users.GET("/me", handler.Profile)
			// add more protected user endpoints here
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Println("listening on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
