package main

import (
	"flowboard-backend-go/internal/database"
	"flowboard-backend-go/internal/middleware"
	"flowboard-backend-go/internal/pages"
	_users "flowboard-backend-go/internal/users"
	"flowboard-backend-go/pkg/config"
	"flowboard-backend-go/pkg/logger"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	logger.Init()
	logger.Log.Infow("Starting FlowBoard API")

	db := database.Connect(cfg)
	db.AutoMigrate(&_users.User{}, &pages.Page{})

	// Users
	userRepo := _users.NewRepository(db)
	userService := _users.NewService(userRepo)
	jwtMgr := middleware.NewJWTManager(cfg.JWTSecret)
	userHandler := _users.NewHandler(userService, jwtMgr)

	// Pages
	pageRepo := pages.NewRepository(db)
	pageService := pages.NewService(pageRepo)
	pageHandler := pages.NewHandler(pageService)

	// Gin
	gin.SetMode(cfg.Mode)
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		auth.POST("/signup", userHandler.Register)
		auth.POST("/login", userHandler.Login)

		usersGroup := api.Group("/users")
		usersGroup.Use(middleware.AuthMiddleware(jwtMgr))
		usersGroup.GET("/me", userHandler.Profile)
	}

	pagesGroup := api.Group("/pages")
	pagesGroup.Use(middleware.AuthMiddleware(jwtMgr))
	pagesGroup.GET("", pageHandler.GetAllPages)
	pagesGroup.POST("", pageHandler.CreatePage)
	pagesGroup.GET("/:id", pageHandler.GetPageByID)
	pagesGroup.PUT("/:id", pageHandler.UpdatePage)
	pagesGroup.DELETE("/:id", pageHandler.DeletePage)

	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Log.Infow("Listening", "port", cfg.Port)
	if err := r.Run(addr); err != nil {
		logger.Log.Fatalw("Server crashed", "error", err)
	}
}
