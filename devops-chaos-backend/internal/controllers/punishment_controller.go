package controllers

import (
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PunishmentController struct {
	punishmentService *services.PunishmentService
}

func NewPunishmentController() *PunishmentController {
	return &PunishmentController{
		punishmentService: services.NewPunishmentService(),
	}
}

// POST /punishments
func (pc *PunishmentController) CreatePunishment(c *gin.Context) {
	var req dto.CreatePunishmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	assignerID, _ := c.Get("userID")
	assignerRole, _ := c.Get("userRole")

	response, err := pc.punishmentService.CreatePunishment(&req, assignerID.(uint), assignerRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to create punishment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Message: "Punishment created successfully",
		Data:    response,
	})
}

// GET /punishments
func (pc *PunishmentController) GetPunishments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	response, err := pc.punishmentService.GetPunishments(page, limit, status, userID.(uint), userRole.(string))
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

// GET /punishments/:id
func (pc *PunishmentController) GetPunishmentByID(c *gin.Context) {
	idParam := c.Param("id")
	punishmentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid punishment ID",
			Error:   err.Error(),
		})
		return
	}

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	response, err := pc.punishmentService.GetPunishmentByID(uint(punishmentID), userID.(uint), userRole.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Message: "Punishment not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Punishment retrieved successfully",
		Data:    response,
	})
}

// PUT /punishments/:id
func (pc *PunishmentController) UpdatePunishment(c *gin.Context) {
	idParam := c.Param("id")
	punishmentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid punishment ID",
			Error:   err.Error(),
		})
		return
	}

	var req dto.UpdatePunishmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	userRole, _ := c.Get("userRole")

	err = pc.punishmentService.UpdatePunishment(uint(punishmentID), &req, userRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to update punishment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Punishment updated successfully",
	})
}

// DELETE /punishments/:id
func (pc *PunishmentController) DeletePunishment(c *gin.Context) {
	idParam := c.Param("id")
	punishmentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid punishment ID",
			Error:   err.Error(),
		})
		return
	}

	userRole, _ := c.Get("userRole")

	err = pc.punishmentService.DeletePunishment(uint(punishmentID), userRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to delete punishment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Punishment deleted successfully",
	})
}

// GET /users/:id/punishments/active
func (pc *PunishmentController) GetActivePunishments(c *gin.Context) {
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

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	punishments, err := pc.punishmentService.GetActivePunishments(uint(targetID), userID.(uint), userRole.(string))
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
		Message: "Active punishments retrieved successfully",
		Data:    punishments,
	})
}
