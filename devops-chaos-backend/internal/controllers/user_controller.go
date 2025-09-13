package controllers

import (
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// GET /users
func (uc *UserController) GetUsers(c *gin.Context) {
	currentUserRole, _ := c.Get("userRole")

	// Parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	role := c.Query("role")

	response, err := uc.userService.GetUsers(page, limit, role, currentUserRole.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Success: false,
			Message: "Access denied",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GET /users/:id
func (uc *UserController) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
		return
	}

	currentUserRole, _ := c.Get("userRole")

	userInfo, err := uc.userService.GetByID(uint(userID), currentUserRole.(string))
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
		Message: "User retrieved successfully",
		Data:    userInfo,
	})
}

// PUT /users/:id
func (uc *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	currentUserRole, _ := c.Get("userRole")

	err = uc.userService.UpdateUser(uint(userID), &req, currentUserRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to update user",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "User updated successfully",
	})
}

// DELETE /users/:id
func (uc *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
		return
	}

	currentUserRole, _ := c.Get("userRole")

	err = uc.userService.DeleteUser(uint(userID), currentUserRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to delete user",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}

// GET /users/stats
func (uc *UserController) GetUserStats(c *gin.Context) {
	currentUserRole, _ := c.Get("userRole")

	stats, err := uc.userService.GetUserStats(currentUserRole.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Success: false,
			Message: "Access denied",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "User statistics retrieved successfully",
		Data:    stats,
	})
}

// POST /users/:id/capture
func (uc *UserController) CaptureNetworkAdmin(c *gin.Context) {
	idParam := c.Param("id")
	targetID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
		return
	}

	currentUserID, _ := c.Get("userID")
	currentUserRole, _ := c.Get("userRole")

	err = uc.userService.CaptureNetworkAdmin(uint(targetID), currentUserID.(uint), currentUserRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to capture network admin",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.CaptureResponse{
		Success:     true,
		Message:     "Network admin captured successfully",
		PointsGiven: 10,
	})
}
