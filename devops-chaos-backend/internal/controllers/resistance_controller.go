package controllers

import (
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResistanceController struct {
	resistanceService *services.ResistanceService
}

func NewResistanceController() *ResistanceController {
	return &ResistanceController{
		resistanceService: services.NewResistanceService(),
	}
}

// GET /resistance
func (rc *ResistanceController) GetResistancePage(c *gin.Context) {
	userRole, _ := c.Get("userRole")

	response, err := rc.resistanceService.GetResistancePage(userRole.(string))
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
		Message: "Resistance page retrieved successfully",
		Data:    response,
	})
}

// POST /resistance/report
func (rc *ResistanceController) ReportSuspiciousActivity(c *gin.Context) {
	var req dto.ReportSuspiciousActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	err := rc.resistanceService.ReportSuspiciousActivity(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to submit report",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Message: "Anonymous report submitted successfully",
	})
}
