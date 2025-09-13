package controllers

import (
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	dashboardService *services.DashboardService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		dashboardService: services.NewDashboardService(),
	}
}

// GET /dashboard/andrei
func (dc *DashboardController) GetAndreiDashboard(c *gin.Context) {
	userRole, _ := c.Get("userRole")

	if userRole != "andrei" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Success: false,
			Message: "Access denied",
			Error:   "Only Andrei can access this dashboard",
		})
		return
	}

	response, err := dc.dashboardService.GetAndreiDashboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve dashboard data",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Andrei dashboard retrieved successfully",
		Data:    response,
	})
}

// GET /dashboard/daemon
// GET /dashboard/daemon (continuación)
func (dc *DashboardController) GetDaemonDashboard(c *gin.Context) {
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	if userRole != "daemon" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Success: false,
			Message: "Access denied",
			Error:   "Only Daemons can access this dashboard",
		})
		return
	}

	response, err := dc.dashboardService.GetDaemonDashboard(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve dashboard data",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Daemon dashboard retrieved successfully",
		Data:    response,
	})
}
