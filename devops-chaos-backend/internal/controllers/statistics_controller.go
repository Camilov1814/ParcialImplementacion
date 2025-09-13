package controllers

import (
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StatisticsController struct {
	statisticsService *services.StatisticsService
}

func NewStatisticsController() *StatisticsController {
	return &StatisticsController{
		statisticsService: services.NewStatisticsService(),
	}
}

// GET /statistics/:user_id
func (sc *StatisticsController) GetUserStatistics(c *gin.Context) {
	idParam := c.Param("user_id")
	targetUserID, err := strconv.ParseUint(idParam, 10, 32)
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

	stats, err := sc.statisticsService.GetUserStatistics(uint(targetUserID), userID.(uint), userRole.(string))
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
		Message: "Statistics retrieved successfully",
		Data:    stats,
	})
}

// GET /statistics/leaderboard
func (sc *StatisticsController) GetLeaderboard(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userRole, _ := c.Get("userRole")

	leaderboard, err := sc.statisticsService.GetLeaderboard(limit, userRole.(string))
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
		Message: "Leaderboard retrieved successfully",
		Data:    leaderboard,
	})
}

// PUT /statistics/:user_id
func (sc *StatisticsController) UpdateStatistics(c *gin.Context) {
	idParam := c.Param("user_id")
	targetUserID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
		return
	}

	var req dto.UpdateStatisticRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	userRole, _ := c.Get("userRole")

	err = sc.statisticsService.UpdateStatistics(uint(targetUserID), &req, userRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to update statistics",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Statistics updated successfully",
	})
}

// POST /statistics/recalculate-rankings
func (sc *StatisticsController) RecalculateRankings(c *gin.Context) {
	userRole, _ := c.Get("userRole")

	err := sc.statisticsService.RecalculateRankings(userRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Failed to recalculate rankings",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Rankings recalculated successfully",
	})
}
