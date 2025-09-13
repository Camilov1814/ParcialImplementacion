package controllers

import (
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	reportService *services.ReportService
}

func NewReportController() *ReportController {
	return &ReportController{
		reportService: services.NewReportService(),
	}
}

// POST /reports
func (rc *ReportController) CreateReport(c *gin.Context) {
	var req dto.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	response, err := rc.reportService.CreateReport(&req, userID.(uint), userRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to create report",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Message: "Report created successfully",
		Data:    response,
	})
}

// GET /reports
func (rc *ReportController) GetReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	reportType := c.Query("type")
	status := c.Query("status")

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	response, err := rc.reportService.GetReports(page, limit, reportType, status, userID.(uint), userRole.(string))
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

// GET /reports/:id
func (rc *ReportController) GetReportByID(c *gin.Context) {
	idParam := c.Param("id")
	reportID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid report ID",
			Error:   err.Error(),
		})
		return
	}

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	response, err := rc.reportService.GetReportByID(uint(reportID), userID.(uint), userRole.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Message: "Report not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Report retrieved successfully",
		Data:    response,
	})
}

// PUT /reports/:id/status
func (rc *ReportController) UpdateReportStatus(c *gin.Context) {
	idParam := c.Param("id")
	reportID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid report ID",
			Error:   err.Error(),
		})
		return
	}

	var req dto.UpdateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	userRole, _ := c.Get("userRole")

	err = rc.reportService.UpdateReportStatus(uint(reportID), &req, userRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to update report status",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Report status updated successfully",
	})
}

// DELETE /reports/:id
func (rc *ReportController) DeleteReport(c *gin.Context) {
	idParam := c.Param("id")
	reportID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid report ID",
			Error:   err.Error(),
		})
		return
	}

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	err = rc.reportService.DeleteReport(uint(reportID), userID.(uint), userRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to delete report",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Report deleted successfully",
	})
}

// GET /reports/recent
func (rc *ReportController) GetRecentReports(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	userRole, _ := c.Get("userRole")

	reports, err := rc.reportService.GetRecentReports(limit, userRole.(string))
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
		Message: "Recent reports retrieved successfully",
		Data:    reports,
	})
}
