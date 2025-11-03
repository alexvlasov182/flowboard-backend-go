package users

import (
	"flowboard-backend-go/internal/auth"
	"flowboard-backend-go/internal/middleware"
	"flowboard-backend-go/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
	jwt     *middleware.JWTManager
}

func NewHandler(service Service, jwt *middleware.JWTManager) *Handler {
	return &Handler{
		service: service,
		jwt:     jwt,
	}
}
func (h *Handler) Register(c *gin.Context) {
	var in auth.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		logger.Log.Warnw("Register input invalid", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(in.Name, in.Email, in.Password)
	if err != nil {
		if err == ErrUserExists {
			logger.Log.Infow("User already exists", "email", in.Email)
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		logger.Log.Errorw("Register failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.jwt.Generate(user.ID)
	if err != nil {
		logger.Log.Errorw("Token generation failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	logger.Log.Infow("User registered successfully", "userID", user.ID, "email", user.Email)
	c.JSON(http.StatusCreated, auth.AuthResponse{
		User:  ToUserResponse(user),
		Token: token,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var in auth.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		logger.Log.Warnw("Login input invalid", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Authenticate(in.Email, in.Password)
	if err != nil {
		logger.Log.Infow("Login failed", "email", in.Email, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := h.jwt.Generate(user.ID)
	if err != nil {
		logger.Log.Errorw("Token generation failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	logger.Log.Infow("User logged in successfully", "userID", user.ID, "email", user.Email)
	c.JSON(http.StatusOK, auth.AuthResponse{
		User:  ToUserResponse(user),
		Token: token,
	})
}

func (h *Handler) Profile(c *gin.Context) {
	uid, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	id := uid.(uint)
	user, err := h.service.GetByID(id)
	if err != nil {
		logger.Log.Errorw("Profile fetch failed", "userID", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	logger.Log.Infow("Profile fetched", "userID", id)
	c.JSON(http.StatusOK, gin.H{"user": ToUserResponse(user)})
}
