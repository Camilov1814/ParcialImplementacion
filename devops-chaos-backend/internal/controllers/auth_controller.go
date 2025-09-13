package controllers

import (
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// POST /auth/login
func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	response, err := ac.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "Authentication failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// POST /auth/register (solo Andrei)
func (ac *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	// Obtener rol del usuario actual del contexto
	currentUserRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "User role not found in context",
		})
		return
	}

	userInfo, err := ac.authService.Register(&req, currentUserRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Registration failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    userInfo,
	})
}

// POST /auth/change-password
func (ac *AuthController) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	// Obtener ID del usuario del contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "User ID not found in context",
		})
		return
	}

	err := ac.authService.ChangePassword(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to change password",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// GET /auth/me
func (ac *AuthController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "User ID not found in context",
		})
		return
	}

	// Obtener información del usuario usando el user service
	userService := services.NewUserService()
	userInfo, err := userService.GetByID(userID.(uint), "")
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    userInfo,
	})
}
