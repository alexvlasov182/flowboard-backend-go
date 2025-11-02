package users

import (
	"flowboard-backend-go/internal/auth"
	"flowboard-backend-go/internal/middleware"
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

// Register route handler
func (h *Handler) Register(c *gin.Context) {
	var in auth.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(in.Name, in.Email, in.Password)
	if err != nil {
		if err == ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.jwt.Generate(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, auth.AuthResponse{
		User:  user,
		Token: token,
	})
}

// Login handler
func (h *Handler) Login(c *gin.Context) {
	var in auth.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Authenticate(in.Email, in.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := h.jwt.Generate(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, auth.AuthResponse{User: user, Token: token})
}

// Example protected route: Get current user profile
func (h *Handler) Profile(c *gin.Context) {
	// middleware already set userID
	uid, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	id := uid.(uint)
	user, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
