package main

import (
	"flowboard-backend-go/internal/database"
	"flowboard-backend-go/internal/middleware"
	"flowboard-backend-go/internal/pages"
	_users "flowboard-backend-go/internal/users"
	"flowboard-backend-go/pkg/logger"
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
		log.Println("No .env file found, using system environment variables")
	}

	logger.Init()
	logger.Log.Infow("Starting FlowBoard API")

	// Check required environment variables
	requiredEnv := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_NAME", "PORT"}
	for _, e := range requiredEnv {
		if os.Getenv(e) == "" {
			logger.Log.Fatalw("Database environment variable not set", "env", e)
		}
	}

	// Connect to database
	db := database.Connect()

	// Auto-migrate database schemas
	if err := db.AutoMigrate(&_users.User{}, &pages.Page{}); err != nil {
		logger.Log.Fatalw("Failed to auto-migrate database schemas", "error", err)
	}

	logger.Log.Infow("Database migrated successfully")

	// create layers users
	repo := _users.NewRepository(db)
	service := _users.NewService(repo)
	jwtMgr := middleware.NewJWTManager()
	handler := _users.NewHandler(service, jwtMgr)

	// pages layers
	pagesRepo := pages.NewRepository(db)
	pagesService := pages.NewService(pagesRepo)
	pageHandler := pages.NewHandler(pagesService)

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
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

	pagesGroup := api.Group("/pages")
	pagesGroup.Use(middleware.AuthMiddleware(jwtMgr))
	{
		pagesGroup.GET("/", pageHandler.GetAllPages)
		pagesGroup.POST("/", pageHandler.CreatePage)
		pagesGroup.GET("/:id", pageHandler.GetPageByID)
		pagesGroup.GET("/user/:userID", pageHandler.GetPageByID)
		pagesGroup.PUT("/:id", pageHandler.UpdatePage)
		pagesGroup.DELETE("/:id", pageHandler.DeletePage)
	}

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	logger.Log.Infow("Listening", "port", port)

	if err := r.Run(addr); err != nil {
		logger.Log.Fatalw("Server crashed", "error", err)
	}
}
